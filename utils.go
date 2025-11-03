package main

import (
	"math/rand"
	"net/http"
	"nyooom/logging"
)

func generateRandomString(length int) string {
	// Charset is URL safe and easy to read
	const charset = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ123456789"

	stringBase := make([]byte, length)
	for i := range stringBase {
		stringBase[i] = charset[rand.Intn(len(charset))]
	}
	return string(stringBase)
}

func httpError(w http.ResponseWriter, message string, code int, err error) {
	http.Error(w, message, code)
	logging.PrintErrStr(message, ": ", err.Error())
}

func httpNewError(w http.ResponseWriter, message string, code int) {
	http.Error(w, message, code)
	logging.PrintErrStr(message)
}
