package main

import (
	"html/template"
	"net/http"
	"nyooom/logging"
)

func epDashboardPage(db AdvancedDB, jwt JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify user is authenticated
		cookie, err := r.Cookie(CookieName)
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		_, ok := jwt.ValidateJWT(cookie.Value)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		// User is authenticated, serve dashboard
		logging.Println("Serving dashboard page")
		tmpl := template.Must(template.ParseFiles("static/dashboard.html"))
		tmpl.Execute(w, nil)
	}
}
