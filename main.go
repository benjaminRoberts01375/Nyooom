package main

import (
	"net/http"
	"nyooom/logging"
	"os"
	"strings"
	"time"
)

func main() {
	// Setup
	logging.ReadConfig()                                        // Setup printing colors
	var db AdvancedDB = SetupDB()                               // Setup database
	jwt := loadJWTSecret(db)                                    // Setup JWT
	devMode := strings.ToLower(os.Getenv("DEV_MODE")) == "true" // Read Dev Mode status

	// Running
	logging.Println("Hello, World")
	setupEndpoints(db, jwt, devMode)

	// Configure server with timeouts
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	server.ListenAndServe()
}

func setupEndpoints(db AdvancedDB, jwt JWTService, devMode bool) {
	// Functional endpoints
	http.HandleFunc("/api/create-link", epCreateLink(db, jwt))
	http.HandleFunc("/api/delete-link", epDeleteLink(db, jwt))
	http.HandleFunc("/api/get-links", epGetLinks(db, jwt))
	http.HandleFunc("/api/login", epLogin(db, jwt))
	http.HandleFunc("/api/create-user", epCreateUser(db, jwt))
	http.HandleFunc("/qr/{id}", epQRCode(db, devMode))
	http.HandleFunc("/{id}", epRedirect(db))

	// UI endpoints
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", noCacheInDevMode(fileServer, devMode)))
	http.HandleFunc("/", epBase(db, jwt))
	http.HandleFunc("/create-account", epCreateUserPage(db))
	http.HandleFunc("/login", epLoginPage(db, jwt))
	http.HandleFunc("/dashboard", epDashboardPage(jwt))
}

// noCacheInDevMode wraps a handler to prevent caching in dev mode
func noCacheInDevMode(h http.Handler, devMode bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if devMode {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}
		h.ServeHTTP(w, r)
	})
}
