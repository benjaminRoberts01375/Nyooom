package main

import (
	"html/template"
	"net/http"
	"nyooom/logging"
)

func epCreateUserPage(db AdvancedDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if a user exists in the database
		userExists, err := db.UserExists(r.Context())
		if err != nil {
			logging.PrintErrStr("Failed to check if user exists: " + err.Error())
			http.Error(w, "Failed to check if user exists: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if userExists {
			logging.Println("User already exists")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		// No user exists, setup account
		logging.Println("Serving HTML")
		tmpl := template.Must(template.ParseFiles("static/create-account.html"))
		tmpl.Execute(w, nil)
	}
}
