package server

import (
	"net/http"
)

func HandleJson(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		h(w, r)
	})
}
