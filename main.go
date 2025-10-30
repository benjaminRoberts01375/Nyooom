package main

import (
	"net/http"
	"nyooom/logging"
)

func main() {
	// Setup
	logging.ReadConfig()          // Setup printing colors
	var db AdvancedDB = SetupDB() // Setup database
	jwt := loadJWTSecret(db)      // Setup JWT

	// Running
	logging.Println("Hello, World")
	setupEndpoints(db, jwt)
	http.ListenAndServe(":8080", nil)
}

func setupEndpoints(db AdvancedDB, jwt JWTService) {
	http.HandleFunc("/api/create-link", epCreateLink(db))
	http.HandleFunc("/api/delete-link", epDeleteLink(db))
	http.HandleFunc("/api/get-links", epGetLinks(db))
	http.HandleFunc("/api/user-exists", epUserExists(db))
	http.HandleFunc("/api/login", epLogin(db, jwt))
	http.HandleFunc("/api/jwt-login", epJWTLogin(jwt))
	http.HandleFunc("/api/create-user", epCreateUser(db))
	http.HandleFunc("/", epRedirect(db))
}
