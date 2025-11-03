package main

import (
	"net/http"
	"nyooom/logging"
	"strings"
	"time"
)

// Clock interface allows for testable time operations
type Clock interface {
	Now() time.Time
}

// RealClock implements Clock using the actual time
type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}

func epRedirect(db AdvancedDB) http.HandlerFunc {
	return epRedirectWithClock(db, RealClock{})
}

func epRedirectWithClock(db AdvancedDB, clock Clock) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.PathValue("id"), "/")
		link, err := db.GetLink(r.Context(), slug)

		if err != nil {
			httpError(w, "Couldn't find the URL you were looking for :(", http.StatusInternalServerError, err)
			return
		}
		err = db.LinkAnalytics(r.Context(), slug, 1, clock.Now())
		if err != nil { // Don't error out, it just sucks
			logging.PrintErrStr("Failed to increment clicks for link " + slug + ": " + err.Error())
		}
		http.Redirect(w, r, "https://"+link.URL, http.StatusMovedPermanently)
	}
}
