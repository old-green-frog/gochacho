package server

import (
	"fmt"
	"gochacho/pkg/conf"
	"gochacho/pkg/migrate"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Server struct {
	config *conf.Config
	db     *sqlx.DB
	router *mux.Router
}

func New() *Server {

	config := conf.New()
	db, _ := sqlx.Open("postgres", config.DatabaseString)
	router := mux.NewRouter()

	return &Server{config, db, router}
}

func (s *Server) Run() {
	migrate.Migrate(s.db)
	s.BuildRotes()
	logrus.Infof("Server is listening on %s:%d", s.config.Host, s.config.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", s.config.Host, s.config.Port), nil)
}
