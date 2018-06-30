package iredmail

import (
	"fmt"
)

func (s *Server) AliasRemove(email string) error {
	isAlias, err := s.regularAliasExists(email)
	if err != nil {
		return err
	}
	if !isAlias {
		return fmt.Errorf("Alias %v does not exist", email)
	}

	tx, err := s.DB.Begin()
	stmt1, err := tx.Prepare("DELETE FROM forwardings WHERE address='" + email + "' and is_list=1")
	_, err = stmt1.Exec()
	stmt2, err := tx.Prepare("DELETE FROM alias WHERE address='" + email + "'")
	_, err = stmt2.Exec()

	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()

	return err
}
