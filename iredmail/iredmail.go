package iredmail

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Server struct {
	DB   *sql.DB
	Data struct {
		Forwardings
	}
}

func (s *Server) Close() {
	s.DB.Close()
}

type queryOptions struct {
	where string
}

func New() (*Server, error) {
	version, err := GetIredMailVersion()
	if err != nil {
		return nil, err
	}
	err = version.Check()
	if err != nil {
		return nil, fmt.Errorf("iredMail version %s is not supported", version)
	}

	if config["port"] == "" {
		config["port"] = "3306"
	}

	if config["host"] == "" {
		config["host"] = "127.0.0.1"
	}

	db, err := sql.Open("mysql", config["user"]+":"+config["password"]+"@tcp("+config["host"]+":"+config["port"]+")/vmail")

	server := &Server{
		DB: db,
	}

	return server, err
}
