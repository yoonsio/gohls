package api

import "net/http"

// index is example handler for index route '/'
func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

// healthz implement healthcheck endpoint that returns 200 status code
func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
