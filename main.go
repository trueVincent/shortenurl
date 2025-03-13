package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/internal/register", RegisterPageHandler).Methods("GET")
	r.HandleFunc("/internal/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/internal/login", LoginHandler).Methods("POST")
	r.HandleFunc("/internal/logout", LogoutHandler).Methods("POST")
	r.HandleFunc("/internal/urlMapping", UrlMappingHandler).Methods("POST")
	r.HandleFunc("/internal/urlMapping/{id}/details", UrlMappingDetailsPageHandler).Methods("GET")
	r.HandleFunc("/{id}", RedirectHandler).Methods("GET")

	return r
}

func main() {
	r := SetupRouter()
	InitializeDatabase()

	// Start server
	port := ":8080"
	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(port, r))
}
