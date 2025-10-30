package main

import (
	"context"
	"net/http"
	"strings"
)

func epRedirect(db AdvancedDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/")
		link, err := db.GetLink(context.Background(), slug)

		if err != nil {
			http.Error(w, "Couldn't find the URL you were looking for :(", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "https://"+link.URL, http.StatusMovedPermanently)
	}
}
