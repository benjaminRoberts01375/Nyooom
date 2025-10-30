package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func epLogin(db AdvancedDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func epCreateUser(db AdvancedDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { // Only allow POST requests
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil {
			// TODO: Handle error
		}
		password := r.Form.Get("password")
		if password == "" {
			// TODO: Handle error
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			// TODO: Handle error
		}
		db.SetUser(r.Context(), hashedPassword)
	}
}

func epUserExists(db AdvancedDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exists, err := db.UserExists(r.Context())
		if err != nil {
			// TODO: Handle error
		}
		if exists {
			// TODO: Redirect to login page
		} else {
			// TODO: Redirect to create account page
		}
	}
}
