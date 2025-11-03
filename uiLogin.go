package main

import (
	"html/template"
	"net/http"
	"nyooom/logging"
)

func epLoginPage(db AdvancedDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if a user exists in the database
		userExists, err := db.UserExists(r.Context())
		if err != nil {
			http.Error(w, "Failed to check if user exists: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if !userExists {
			logging.Println("User does not exist")
			http.Redirect(w, r, "/create-account", http.StatusFound)
			return
		}
		// User exists, setup account
		logging.Println("Serving user login HTML")
		tmpl := template.Must(template.ParseFiles("static/login.html"))
		tmpl.Execute(w, nil)
	}
}
