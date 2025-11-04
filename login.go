package main

import (
	"errors"
	"net/http"
	"nyooom/logging"

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
		err = jwtService.setJWT(w)
		if err != nil {
			httpError(w, "Failed to generate JWT", http.StatusInternalServerError, err)
			return
		}
	}
}

func epCreateUser(db AdvancedDB, jwt JWTService) http.HandlerFunc {
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
		err = jwtService.setJWT(w)
		if err != nil {
			httpError(w, "Failed to generate JWT", http.StatusInternalServerError, err)
			return
		}
		http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
	}
}

func (s *JWTService) ValidatePassword(passwordAttempt string, realPasswordHash []byte) bool {
	return bcrypt.CompareHashAndPassword(realPasswordHash, []byte(passwordAttempt)) == nil
}
