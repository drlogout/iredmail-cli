package integrationTest

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbConnectionString = "vmail:sx4fDttWdWNbiBPsGxhbbxic2MmmGsmJ@tcp(127.0.0.1:8806)/vmail"
)

var dbTables = []string{
	"alias",
	"domain",
	"forwardings",
	"mailbox",
}

func setupDB() error {
	db, err := sql.Open("mysql", dbConnectionString)
	if err != nil {
		return err
	}

	for _, table := range dbTables {
		_, err := db.Exec("DELETE FROM " + table)
		if err != nil {
			return err
		}
	}

	return nil
}
