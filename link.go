package main

import (
	"errors"
	"strings"
)

type Link struct {
	Slug   string
	URL    string
	ID     string
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
		ID:     generateRandomString(10),
		Clicks: 0,
	}, nil
}
