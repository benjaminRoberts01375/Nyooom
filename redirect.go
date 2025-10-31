package main

import (
	"net/http"
	"nyooom/logging"
	"strings"
)

func epRedirect(db AdvancedDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.PathValue("id"), "/")
		link, err := db.GetLink(r.Context(), slug)

		if err != nil {
			http.Error(w, "Couldn't find the URL you were looking for :(", http.StatusInternalServerError)
			return
		}
		err = db.IncrementLinkClicks(r.Context(), slug, 1)
		if err != nil {
			logging.PrintErrStr("Failed to increment clicks for link " + slug + ": " + err.Error())
		}
		http.Redirect(w, r, "https://"+link.URL, http.StatusMovedPermanently)
	}
}
