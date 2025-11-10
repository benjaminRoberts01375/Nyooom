package main

import (
	"html/template"
	"net/http"
	"nyooom/logging"
)

func epDashboardPage(jwt JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify user is authenticated
		err := jwt.ReadAndValidateJWT(r)
		if err != nil {
			logging.Println("JWT is invalid, redirecting to login page")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		// User is authenticated, serve dashboard
		logging.Println("Serving dashboard page")
		tmpl := template.Must(template.ParseFiles("static/dashboard.html"))
		tmpl.Execute(w, nil)
	}
}
