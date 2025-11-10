package main

import (
	"net/http"
	"nyooom/logging"
	"time"
)

func main() {
	// Setup
	logging.ReadConfig()          // Setup printing colors
	var db AdvancedDB = SetupDB() // Setup database
	jwt := loadJWTSecret(db)      // Setup JWT

	// Running
	logging.Println("Hello, World")
	setupEndpoints(db, jwt)

	// Configure server with timeouts
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	server.ListenAndServe()
}

func setupEndpoints(db AdvancedDB, jwt JWTService) {
	// Functional endpoints
	http.HandleFunc("/api/create-link", epCreateLink(db, jwt))
	http.HandleFunc("/api/delete-link", epDeleteLink(db, jwt))
	http.HandleFunc("/api/get-links", epGetLinks(db, jwt))
	http.HandleFunc("/api/login", epLogin(db, jwt))
	http.HandleFunc("/api/create-user", epCreateUser(db, jwt))
	http.HandleFunc("/qr/{id}", epQRCode(db))
	http.HandleFunc("/{id}", epRedirect(db))

	// UI endpoints
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", epBase(db, jwt))
	http.HandleFunc("/create-account", epCreateUserPage(db))
	http.HandleFunc("/login", epLoginPage(db, jwt))
	http.HandleFunc("/dashboard", epDashboardPage(jwt))
}
