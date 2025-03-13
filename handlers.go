package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

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
			urlList, err := ListUrlMapping(userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
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

	_, err := Register(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := Login(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	userID := session.Values["userID"].(uint)

	_, err = CreateUrlMapping(userID, originURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	urlMapping, err := GetUrlMapping(urlMappingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	urlMappingActionRecord, err := GetUrlMappingActionRecord(urlMappingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	urlMapping, err := Redirect(urlMappingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, urlMapping.OriginURL, http.StatusFound)
}
