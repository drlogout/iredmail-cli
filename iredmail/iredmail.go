package iredmail

import (
	"database/sql"
	"fmt"

	"github.com/spf13/viper"

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

	db, err := sql.Open("mysql", viper.GetString("dbuser")+":"+viper.GetString("dbpassword")+"@tcp("+viper.GetString("dbhost")+":"+viper.GetString("dbport")+")/vmail")

	server := &Server{
		DB: db,
	}

	return server, err
}
