package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data map[string]interface{}) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	session, err := store.Get(r, "shortenurl_session")
	login := session.Values["login"]
	if login == nil {
		login = false
	}
	data["Login"] = login.(bool)
	data["Username"] = session.Values["username"].(string)

	if err := t.Execute(w, data); err != nil {
		log.Printf("Template execute error: %v", err)
		http.Error(w, "Template execute error", http.StatusInternalServerError)
		return
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	session, err := store.Get(r, "shortenurl_session")
	if err != nil {
		log.Printf("Session error: %v", err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	if login, ok := session.Values["login"].(bool); ok && login {
		userID, ok := session.Values["userID"].(uint)
		if ok {
			var urlList []URLMapping
			if err := DB.Where("user_id = ?", userID).Find(&urlList).Error; err != nil {
				log.Printf("Database error: %v", err)
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}

			data["UrlList"] = urlList
		}
	} else {
		session.Values["login"] = false
		session.Values["username"] = ""
		session.Save(r, w)
	}
	RenderTemplate(w, r, "index.html", data)
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "register.html", make(map[string]interface{}))
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))
	if username == "" || password == "" {
		http.Error(w, "Missing username or password", http.StatusBadRequest)
		return
	}

	var existUserCnt int64
	if err := DB.Model(&User{}).Where("username = ?", username).Count(&existUserCnt).Error; err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if existUserCnt != 0 {
		data := map[string]interface{}{}
		data["Error"] = "username already used"
		RenderTemplate(w, r, "register.html", data)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	newUser := User{Username: username, Password: string(hashedPassword)}
	result := DB.Create(&newUser)
	if result.Error != nil {
		log.Printf("Error creating user")
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	var user User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		log.Printf("Invalid username or password %v", err)
		data := make(map[string]interface{})
		data["Error"] = "Invalid username or password"
		RenderTemplate(w, r, "index.html", data)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("Invalid username or password %v", err)
		data := make(map[string]interface{})
		data["Error"] = "Invalid username or password"
		RenderTemplate(w, r, "index.html", data)
		return
	}

	session, _ := store.Get(r, "shortenurl_session")
	session.Values["login"] = true
	session.Values["username"] = username
	session.Values["userID"] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "shortenurl_session")
	session.Values["login"] = false
	session.Values["username"] = ""
	session.Values["userID"] = ""
	session.Save(r, w)
	
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UrlMappingHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "shortenurl_session")
	if err != nil {
		log.Printf("Session error: %v", err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	if session.Values["login"] == false {
		http.Error(w, "Shorten URL after login", http.StatusUnauthorized)
		return
	}

	originURL := strings.TrimSpace(r.FormValue("originURL"))

	var randomID string
	for randomID == "" {
		randomID = randomString(6)
		var existUrlMappingCnt int64
		if err := DB.Model(&URLMapping{}).Where("ID = ?", randomID).Count(&existUrlMappingCnt).Error; err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if existUrlMappingCnt != 0 {
			randomID = ""
		}
	}

	newURLMapping := URLMapping{ID: randomID, OriginURL: originURL, UserID: session.Values["userID"].(uint)}
	if err := DB.Create(&newURLMapping).Error; err != nil {
		log.Printf("Error creating url mapping: %v", err)
		http.Error(w, "Error creating url mapping", http.StatusInternalServerError)
		return
	}
	newURLMappingActionRecord := URLMappingActionRecord{URLMapping: newURLMapping, ClickCount: 0}
	if err := DB.Create(&newURLMappingActionRecord).Error; err != nil {
		log.Printf("Error creating url mapping action record: %v", err)
		http.Error(w, "Error creating url mapping action record", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UrlMappingDetailsPageHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "shortenurl_session")
	if err != nil {
		log.Printf("Session error: %v", err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	if login, ok := session.Values["login"].(bool); !ok || !login {
		http.Error(w, "Shorten URL after login", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	urlMappingID := vars["id"]
	var urlMapping URLMapping
	if err := DB.Where("ID = ?", urlMappingID).Take(&urlMapping).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		log.Printf("Error quering URLMapping: %v", err)
		http.Error(w, "Error quering URLMapping", http.StatusInternalServerError)
		return
	}
	var urlMappingActionRecord URLMappingActionRecord
	if err := DB.Where("url_mapping_id = ?", urlMappingID).Take(&urlMappingActionRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "URL action record not found", http.StatusNotFound)
			return
		}
		log.Printf("Error quering URLMappingActionRecord: %v", err)
		http.Error(w, "Error quering URLMappingActionRecord", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{}
	data["UrlMapping"] = urlMapping
	data["UrlMappingActionRecord"] = urlMappingActionRecord

	RenderTemplate(w, r, "urlMapping-details.html", data)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlMappingID := vars["id"]

	urlMappingChan := make(chan URLMapping, 1)
	urlMappingActionRecordChan := make(chan URLMappingActionRecord, 1)
	errChan := make(chan error, 2)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		var urlMapping URLMapping
		if err := DB.Where("ID =?", urlMappingID).Take(&urlMapping).Error; err != nil {
			errChan <- err
			return
		}
		urlMappingChan <- urlMapping
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var urlMappingActionRecord URLMappingActionRecord
		if err := DB.Where("url_mapping_id = ?", urlMappingID).Take(&urlMappingActionRecord).Error; err != nil {
			errChan <- err
			return
		}
		urlMappingActionRecord.ClickCount += 1
		urlMappingActionRecord.LastAccess = time.Now()
		if err := DB.Save(&urlMappingActionRecord).Error; err != nil {
			errChan <- err
			return
		}
		urlMappingActionRecordChan <- urlMappingActionRecord
	}()

	wg.Wait()
	close(urlMappingChan)
	close(urlMappingActionRecordChan)
	close(errChan)

	for err := range errChan {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "URL not found", http.StatusNotFound)
		} else {
			log.Printf("Database error %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	urlMapping := <-urlMappingChan
	http.Redirect(w, r, urlMapping.OriginURL, http.StatusFound)
}
