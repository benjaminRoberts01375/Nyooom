package main

import (
	"net/http"
	"nyooom/logging"
)

func main() {
	logging.ReadConfig() // Setup printing colors
	var db AdvancedDB = SetupDB()
	logging.Println("Hello, World")
	setupEndpoints(db)
	http.ListenAndServe(":8080", nil)
}

func setupEndpoints(db AdvancedDB) {
	http.HandleFunc("/create-link", epCreateLink(db))
	http.HandleFunc("/delete-link", epDeleteLink(db))
	http.HandleFunc("/get-links", epGetLinks(db))
}
