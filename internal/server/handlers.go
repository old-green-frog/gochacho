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
	s.router.PathPrefix("/favicon").Handler(http.FileServer(http.Dir("static/favicon")))

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
	reports.HandleFunc("/{option}", s.ReportsHandler).Methods("GET", "POST")
	reports.HandleFunc("/{option}/create", CreateReportHandler).Methods("GET", "CREATE")

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
				Name:  "role",
				Value: role,
				Path:  "/",
			}
			http.SetCookie(w, cook)
			http.Redirect(w, r, fmt.Sprintf("/reports/%s", option), http.StatusFound)
		}
	}
}

func (s *Server) ReportsHandler(w http.ResponseWriter, r *http.Request) {

	key := mux.Vars(r)
	roleCookie, err := r.Cookie("role")
	if err != nil {
		logrus.Error(err)
		http.HandlerFunc(ServNotFoundHandler).ServeHTTP(w, r)
		return
	}

	option := key["option"]
	role := roleCookie.Value

	switch r.Method {
	case "GET":
		{
			context := ROLES_OPTIONS[role]
			var element TemplateOption

			for _, el := range context {
				if el.Key == option {
					element = el
				}
			}
			if element.Name == "" {
				logrus.Error(err)
				http.HandlerFunc(ServNotFoundHandler).ServeHTTP(w, r)
				return
			}

			IsNumber := false
			if option == "defectlist" {
				IsNumber = true
			}

			tempateContext := struct {
				Title    string
				IsDate   bool
				Results  bool
				IsNumber bool
				Option   string
				IsCreate bool
			}{element.Name, element.Date, false, IsNumber, option, element.Create}

			tmpl, _ := template.ParseFiles("static/view.html")
			tmpl.Execute(w, tempateContext)
			return
		}

	case "POST":
		{
			r.ParseForm()

			var formData []string

			for key, value := range r.Form {
				if key != "sign-button" && key != "create-button" {
					formData = append(formData, value[0])
				}
			}

			context := ROLES_OPTIONS[role]
			var element TemplateOption

			for _, el := range context {
				if el.Key == option {
					element = el
				}
			}

			titles, data := s.getReportData(option, formData)

			IsNumber := false
			if option == "defectlist" {
				IsNumber = true
			}

			tempateContext := struct {
				Title      string
				IsDate     bool
				Titles     []string
				Data       [][]string
				Results    bool
				FirstDate  string
				SecondDate string
				IsNumber   bool
				Option     string
				IsCreate   bool
			}{
				element.Name,
				element.Date,
				titles,
				data,
				true,
				formData[0],
				"",
				IsNumber,
				option,
				element.Create,
			}

			if len(formData) > 1 {
				tempateContext.SecondDate = formData[1]
			}
			if len(data) == 0 && len(titles) != 0 {
				tempateContext.Results = false
			}

			tmpl, _ := template.ParseFiles("static/view.html")
			tmpl.Execute(w, tempateContext)
		}
	}
}

func CreateReportHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	option := vars["option"]

	tempateContext := struct {
		Titles []string
	}{
		TITLES[option],
	}

	tmpl, _ := template.ParseFiles("static/create.html")
	tmpl.Execute(w, tempateContext)
}
