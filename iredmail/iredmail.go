package iredmail

import (
	"bufio"
	"database/sql"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Server struct {
	DB *sql.DB
}

type queryOptions struct {
	where string
}

const (
	myCnfName = ".my.cnf-vmail"
)

func New() (*Server, error) {
	dbConfig, _ := parseMyCnf()
	db, err := sql.Open("mysql", dbConfig["user"]+":"+dbConfig["password"]+"@tcp("+dbConfig["host"]+":"+dbConfig["port"]+")/vmail")

	server := &Server{
		DB: db,
	}

	return server, err
}

func parseMyCnf() (map[string]string, error) {
	dbConfig := map[string]string{}
	usr, err := user.Current()
	if err != nil {
		return dbConfig, err
	}

	myCnf, err := os.Open(filepath.Join(usr.HomeDir, myCnfName))
	if err != nil {
		return dbConfig, err
	}
	defer myCnf.Close()

	scanner := bufio.NewScanner(myCnf)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, "=")
		if len(splitLine) > 1 {
			dbConfig[splitLine[0]] = strings.Trim(splitLine[1], "\"")
		}
	}

	return dbConfig, nil
}
