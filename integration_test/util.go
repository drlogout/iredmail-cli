package test

import "database/sql"

var dbTables = []string{
	"alias",
	"domain",
	"forwardings",
	"mailbox",
}

func setupDB() error {
	db, err := sql.Open("mysql", "vmail:sx4fDttWdWNbiBPsGxhbbxic2MmmGsmJ@tcp(127.0.0.1:8806)/vmail")
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
