package main

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// nopCloser wraps an io.Writer to make it an io.WriteCloser
type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func epQRCode(db AdvancedDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.PathValue("id"), "/")

		// Verify the link exists
		_, err := db.GetLink(r.Context(), slug)
		if err != nil {
			httpError(w, "Couldn't find the URL you were looking for :(", http.StatusNotFound, err)
			return
		}

		// Generate the full shortened URL
		host := r.Host
		if host == "" {
			host = "localhost:8080"
		}
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		shortURL := scheme + "://" + host + "/" + slug

		// Generate QR code
		qrc, err := qrcode.New(shortURL)
		if err != nil {
			httpError(w, "Failed to generate QR code", http.StatusInternalServerError, err)
			return
		}

		// Write QR code to buffer
		buf := new(bytes.Buffer)
		w2 := standard.NewWithWriter(nopCloser{buf}, standard.WithQRWidth(21))
		if err = qrc.Save(w2); err != nil {
			httpError(w, "Failed to save QR code", http.StatusInternalServerError, err)
			return
		}

		// Send the QR code as PNG
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 1 day
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	}
}
