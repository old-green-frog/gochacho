package server

import (
	"database/sql"
	"fmt"
	"gochacho/pkg/conf"
	"gochacho/pkg/migrate"

	_ "github.com/lib/pq"
)

type Server struct {
	config *conf.Config
	db     *sql.DB
}

func New() *Server {

	config := conf.New()
	base, _ := sql.Open("postgres", config.DatabaseString)

	return &Server{
		config: config,
		db:     base,
	}
}

func (s *Server) Run() {
	migrate.Migrate(s.db)
	fmt.Println(s.config.DatabaseString)
}
