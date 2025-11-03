package main

import (
	"net/http"
)

func epBase(db AdvancedDB, jwt JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if a user exists in the database
		userExists, err := db.UserExists(r.Context())
		if err != nil {
			httpError(w, "Failed to check if user exists: ", http.StatusInternalServerError, err)
			return
		}
		if !userExists {
			http.Redirect(w, r, "/create-account", http.StatusTemporaryRedirect)
			return
		}
		// Check if the user is authenticated
		cookie, err := r.Cookie(CookieName)
		if err != nil || cookie.Value == "" { // No or bad cookie
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		_, ok := jwt.ValidateJWT(cookie.Value)
		if !ok { // Invalid JWT
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		// All is good, send to the dashboard
		http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
	}
}
