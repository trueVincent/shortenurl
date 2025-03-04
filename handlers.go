package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/gorilla/mux"
)

// TODO: secret key
var store = sessions.NewCookieStore([]byte("your-secret-key"))

func RenderTemplate(w http.ResponseWriter, tmpl string, data ...interface{}) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: index.html has title, login btn, register link, shorten btn, url list
	data := map[string]interface{}{}
	
	// TODO: session_name
	session, err := store.Get(r, "session_name")
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

			data["urlList"] = urlList
		}
	}
	RenderTemplate(w, "index.html", data)
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: register.html has register btn
	RenderTemplate(w, "register.html")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Handle registration logic (e.g., save user to DB), the password need to be handled carefully
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Handle login logic (e.g., authenticate user), update session
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UrlMappingHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Handle URL shortening logic
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UrlMappingDetailsPageHandler(w http.ResponseWriter, r *http.Request) {
	// urlMapping-details.html has url, count, last_access
	RenderTemplate(w, "urlMapping-details.html")
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// TODO: Lookup original URL from DB
	// TODO: async update count and last_access
	originalURL := "https://example.com/original-url"

	log.Printf("Redirecting %s to %s\n", id, originalURL)
	http.Redirect(w, r, originalURL, http.StatusFound)
}
