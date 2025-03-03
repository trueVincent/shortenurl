package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func RenderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "index.html")
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "register.html")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Handle registration logic (e.g., save user to DB)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Handle login logic (e.g., authenticate user)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UrlMappingHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Handle URL shortening logic
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UrlMappingDetailsHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "urlMapping-details.html")
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// TODO: Lookup original URL from DB
	originalURL := "https://example.com/original-url"

	log.Printf("Redirecting %s to %s\n", id, originalURL)
	http.Redirect(w, r, originalURL, http.StatusFound)
}
