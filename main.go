package main

import (
	"net/http"
	"nyooom/logging"
)

func main() {
	logging.ReadConfig()          // Setup printing colors
	var db AdvancedDB = SetupDB() // Setup database
	logging.Println("Hello, World")
	setupEndpoints(db)
	http.ListenAndServe(":8080", nil)
}

func setupEndpoints(db AdvancedDB) {
	http.HandleFunc("/api/create-link", epCreateLink(db))
	http.HandleFunc("/api/delete-link", epDeleteLink(db))
	http.HandleFunc("/api/get-links", epGetLinks(db))
}
