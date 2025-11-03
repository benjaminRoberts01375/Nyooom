package main

import (
	"errors"
	"net/http"
	"nyooom/logging"
	"strconv"
	"strings"
)

type Link struct {
	Slug   string
	URL    string
	Clicks int
}

func newLink(slug string, url string) (Link, error) {
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

func epCreateLink(db AdvancedDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		link, err := newLink(r.URL.Query().Get("slug"), r.URL.Query().Get("url"))
		if err != nil {
			httpError(w, "Failed to create link \""+link.Slug+".\"", http.StatusInternalServerError, err)
			return
		}
		err = db.SetLink(r.Context(), link)
		if err != nil {
			httpError(w, "Failed to create link \""+link.Slug+".\" in database", http.StatusInternalServerError, err)
			return
		}
		// TODO: Handle success
		logging.Println("Created link \"" + link.Slug + "\"")
	}
}

func epDeleteLink(db AdvancedDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		linkSlug := r.URL.Query().Get("slug")
		err := db.DeleteLink(r.Context(), linkSlug)
		if err != nil {
			httpError(w, "Failed to delete link \""+linkSlug+".\" ", http.StatusInternalServerError, err)
			return
		}
		logging.Println("Deleted link \"" + linkSlug + "\"")
		// TODO: Handle success
	}
}

func epGetLinks(db AdvancedDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		links, err := db.GetLinks(r.Context())
		if err != nil {
			httpError(w, "Failed to get links", http.StatusInternalServerError, err)
			return
		}
		// TODO Handle success
		for _, link := range links {
			logging.Println(link)
		}
	}
}
