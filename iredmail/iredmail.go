package iredmail

import (
	"database/sql"

	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
)

type Server struct {
	DB *sql.DB
}

func (s *Server) Close() {
	s.DB.Close()
}

type queryOptions struct {
	where string
}

func New() (*Server, error) {
	db, err := sql.Open("mysql", viper.GetString("dbuser")+":"+viper.GetString("dbpassword")+"@tcp("+viper.GetString("dbhost")+":"+viper.GetString("dbport")+")/vmail")

	server := &Server{
		DB: db,
	}

	return server, err
}
