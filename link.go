package main

import (
	"errors"
	"html/template"
	"net/http"
	"nyooom/logging"
	"strconv"
	"strings"
	"time"
)

// Helper function to render link cards template
func renderLinkCards(w http.ResponseWriter, links []Link) error {
	tmpl, err := template.ParseFiles("static/link-cards.html")
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/html")
	return tmpl.Execute(w, links)
}

type Link struct {
	Slug      string
	URL       string
	Clicks    int
	LastClick *time.Time
}

func newLink(slug string, url string) (Link, error) {
	logging.Println("Slug: ", slug)
	logging.Println("URL:  ", url)
	url, _ = strings.CutPrefix(url, "https://")
	url, _ = strings.CutPrefix(url, "http://")

	// Check if slug and URL are set
	if len(slug) < 3 || // Slug must be at least 3 characters
		strings.Contains(slug, " ") || // Slug cannot contain spaces
		strings.Contains(url, " ") || // URL cannot contain spaces
		!strings.Contains(url, ".") || // URL must contain a dot
		len(url) < 5 { // URL must be at least 5 characters
		return Link{}, errors.New("Invalid slug or URL: " + slug + " | " + url)
	}

	return Link{
		Slug:   slug,
		URL:    url,
		Clicks: 0,
	}, nil
}

func (link Link) String() string {
	return link.Slug + " -> " + link.URL + " has " + strconv.Itoa(link.Clicks) + " clicks"
}

func epCreateLink(db AdvancedDB, jwt JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify user is authenticated
		cookie, err := r.Cookie(CookieName)
		if err != nil || cookie.Value == "" {
			httpError(w, "Unauthorized", http.StatusUnauthorized, err)
			return
		}
		_, ok := jwt.ValidateJWT(cookie.Value)
		if !ok {
			httpError(w, "Unauthorized", http.StatusUnauthorized, errors.New("invalid JWT"))
			return
		}

		link, err := newLink(r.FormValue("slug"), r.FormValue("url"))
		if err != nil {
			httpError(w, "Failed to create link \""+link.Slug+".\"", http.StatusBadRequest, err)
			return
		}
		err = db.SetLink(r.Context(), link)
		if err != nil {
			httpError(w, "Failed to create link \""+link.Slug+".\" in database", http.StatusInternalServerError, err)
			return
		}
		logging.Println("Created link \"" + link.Slug + "\"")

		// Return the updated links list for HTMX
		links, err := db.GetLinks(r.Context())
		if err != nil {
			httpError(w, "Failed to get links", http.StatusInternalServerError, err)
			return
		}

		// Render links using template
		err = renderLinkCards(w, links)
		if err != nil {
			httpError(w, "Failed to render links", http.StatusInternalServerError, err)
		}
	}
}

func epDeleteLink(db AdvancedDB, jwt JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify user is authenticated
		cookie, err := r.Cookie(CookieName)
		if err != nil || cookie.Value == "" {
			httpError(w, "Unauthorized", http.StatusUnauthorized, err)
			return
		}
		_, ok := jwt.ValidateJWT(cookie.Value)
		if !ok {
			httpError(w, "Unauthorized", http.StatusUnauthorized, errors.New("invalid JWT"))
			return
		}

		linkSlug := r.URL.Query().Get("slug")
		err = db.DeleteLink(r.Context(), linkSlug)
		if err != nil {
			httpError(w, "Failed to delete link \""+linkSlug+".\" ", http.StatusInternalServerError, err)
			return
		}
		logging.Println("Deleted link \"" + linkSlug + "\"")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Link deleted successfully"))
	}
}

func epGetLinks(db AdvancedDB, jwt JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify user is authenticated
		cookie, err := r.Cookie(CookieName)
		if err != nil || cookie.Value == "" {
			httpError(w, "Unauthorized", http.StatusUnauthorized, err)
			return
		}
		_, ok := jwt.ValidateJWT(cookie.Value)
		if !ok {
			httpError(w, "Unauthorized", http.StatusUnauthorized, errors.New("invalid JWT"))
			return
		}

		links, err := db.GetLinks(r.Context())
		if err != nil {
			httpError(w, "Failed to get links", http.StatusInternalServerError, err)
			return
		}

		// Render links using template
		err = renderLinkCards(w, links)
		if err != nil {
			httpError(w, "Failed to render links", http.StatusInternalServerError, err)
		}
	}
}
