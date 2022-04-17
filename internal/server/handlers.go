package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

var ROLES = map[string]string{
	"mech":    "Механик",
	"main_ch": "Руководитель",
	"caseer":  "Кассир-контролер",
}

const (
	ROLE_QUERY = "SELECT id FROM adminroles WHERE val='%s'"
	AUTH_QUERY = "SELECT * FROM admins WHERE serial_num='%s' AND role_id=(" + ROLE_QUERY + ")"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		logrus.Infof("%s\t%s\n", r.RequestURI, r.Method)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (s *Server) BuildRotes() {

	// static files
	statics := s.router.PathPrefix("/static").Subrouter()
	statics.PathPrefix("/css").Handler(http.StripPrefix("/static/css/", http.FileServer(http.Dir("static/css"))))
	statics.PathPrefix("/js").Handler(http.StripPrefix("/static/js/", http.FileServer(http.Dir("static/js"))))
	statics.PathPrefix("/images").Handler(http.StripPrefix("/static/images/", http.FileServer(http.Dir("static/images"))))

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	}).Methods("GET")

	s.router.HandleFunc("/usernf", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/usernf.html")
	}).Methods("GET")

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		number := r.Form.Get("number")
		role := ""
		rkey := ""

		for key := range r.Form {
			if key != "sign-button" && key != "number" {
				rkey = key
				role = ROLES[key]
			}
		}

		if auth(number, role, s.db) {
			http.Redirect(w, r, rkey, http.StatusMovedPermanently)
		} else {
			http.Redirect(w, r, "usernf", http.StatusMovedPermanently)
		}
	}).Methods("POST")

	s.router.NotFoundHandler = http.HandlerFunc(ServNotFoundHandler)
	s.router.Use(loggingMiddleware)

	http.Handle("/", s.router)
}

func auth(num, role string, db *sql.DB) bool {
	if len(num) > 15 {
		return false
	}
	err := db.QueryRow(fmt.Sprintf(AUTH_QUERY, num, role)).Scan()
	if err != nil {
		logrus.Warn(err)
		return false
	}
	return true
}

func ServNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/404.html")
}
