package server

import "net/http"

func Handle404(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not found!", http.StatusNotFound)
}
