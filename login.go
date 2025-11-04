package main

import (
	"errors"
	"net/http"
	"nyooom/logging"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func epLogin(db AdvancedDB, jwtService JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { // Only allow POST requests
			httpNewError(w, "Method not allowed for password login", http.StatusMethodNotAllowed)
			return
		}

		exists, err := db.UserExists(r.Context())
		if err != nil {
			httpError(w, "Could not check if user already exists", http.StatusInternalServerError, err)
			return
		}
		if !exists {
			logging.Println("No users exist")
			http.Redirect(w, r, "/create-account", http.StatusTemporaryRedirect)
			return
		}

		err = r.ParseForm()
		if err != nil {
			httpError(w, "Could not read password", http.StatusBadRequest, err)
			return
		}
		password := r.Form.Get("password")
		if password == "" {
			httpError(w, "Could not read password", http.StatusBadRequest, errors.New("Password field is blank or missing"))
			return
		}

		realHash, err := db.GetUser(r.Context())
		if err != nil {
			httpError(w, "Could not validate password", http.StatusInternalServerError, err)
			return
		}
		if !jwtService.ValidatePassword(password, []byte(realHash)) {
			httpError(w, "Incorrect password", http.StatusForbidden, err)
			return
		}
		// Generate a new JWT
		token, err := jwtService.GenerateJWT(LoginDuration)
		if err != nil {
			httpError(w, "Failed to generate JWT", http.StatusInternalServerError, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     CookieName,
			Value:    token,
			HttpOnly: false,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(LoginDuration),
			Path:     "/",
		})

	}
}

func epJWTLogin(jwtService JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { // Only allow POST requests
			httpNewError(w, "Method not allowed for JWT login", http.StatusMethodNotAllowed)
			return
		}
		err := r.ParseForm()
		if err != nil {
			httpNewError(w, "Failed to parse JWT form", http.StatusBadRequest)
			return
		}
		password := r.Form.Get("jwt")
		if password == "" {
			httpNewError(w, "Password is required for login", http.StatusBadRequest)
			return
		}
		_, valid := jwtService.ValidateJWT(password)
		if !valid {
			httpNewError(w, "Bad JWT", http.StatusNotAcceptable)
			return
		}
		http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
	}
}

func epCreateUser(db AdvancedDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { // Only allow POST requests
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		exists, err := db.UserExists(r.Context())
		if exists {
			httpNewError(w, "User already exists", http.StatusForbidden)
			return
		} else if err != nil {
			httpError(w, "Failed to check if user exists", http.StatusInternalServerError, err)
			return
		}

		err = r.ParseForm()
		if err != nil {
			httpError(w, "Failed to parse form", http.StatusBadRequest, err)
			return
		}
		password := r.Form.Get("password")
		if password == "" {
			httpNewError(w, "Password is required", http.StatusBadRequest)
			return
		}
		if len(password) < 8 {
			httpNewError(w, "Password must be at least 8 characters", http.StatusBadRequest)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			httpError(w, "Failed to hash new password", http.StatusInternalServerError, err)
			return
		}
		err = db.SetUser(r.Context(), hashedPassword)
		if err != nil {
			httpError(w, "Failed to create user", http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusCreated)

		// TODO: Redirect to dashboard
	}
}

func (s *JWTService) ValidatePassword(passwordAttempt string, realPasswordHash []byte) bool {
	return bcrypt.CompareHashAndPassword(realPasswordHash, []byte(passwordAttempt)) == nil
}
