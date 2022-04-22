package server

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		role := ""

		for key := range r.Form {
			if key != "sign-button" {
				role = key
			}
		}

		http.Redirect(w, r, "/info/"+role, http.StatusMovedPermanently)
	}).Methods("POST")

	s.router.HandleFunc("/info/{role}", InfoHandler).Methods("GET", "POST")

	reports := s.router.PathPrefix("/reports").Subrouter()
	reports.HandleFunc("/{option}", ReportsHandler).Methods("GET", "POST")
	reports.HandleFunc("/{option}/create", CreateReportHandler).Methods("CREATE")

	s.router.NotFoundHandler = http.HandlerFunc(ServNotFoundHandler)
	s.router.Use(loggingMiddleware)

	http.Handle("/", s.router)
}

func ServNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/404.html")
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {

	key := mux.Vars(r)
	role := key["role"]

	if x := r.URL.Query().Get("option"); x == "" && r.Method == "GET" {
		context := TemplateData{ROLES_OPTIONS[role]}
		tmpl, _ := template.ParseFiles("static/info.html")
		tmpl.Execute(w, context)
		return
	}

	if r.Method == "POST" {
		context := TemplateData{ROLES_OPTIONS[role]}
		r.ParseForm()

		option := ""
		for key := range r.Form {
			if key != "sign-button" {
				option = key
			}
		}
		var element TemplateOption
		for _, el := range context.Options {
			logrus.Info(el)
			if el.Key == option && option != "" {
				element = el
			}
		}
		if element.Key != "" {
			cook := &http.Cookie{
				Name:  "is_date",
				Value: fmt.Sprint(element.Date),
				Path:  "/",
			}
			http.SetCookie(w, cook)
			http.Redirect(w, r, fmt.Sprintf("/reports/%s", option), http.StatusFound)
		}
	}
}

func ReportsHandler(w http.ResponseWriter, r *http.Request) {

}

func CreateReportHandler(w http.ResponseWriter, r *http.Request) {

}
