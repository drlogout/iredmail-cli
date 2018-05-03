package iredmail

import "fmt"

func (s *Server) AliasRemove(email string) error {
	isAlias, err := s.isAlias(email)
	if err != nil {
		return err
	}

	if isAlias {
		_, err = s.DB.Exec(`
			DELETE FROM forwardings WHERE address='` + email + `' and is_list=1
		`)
		if err != nil {
			return err
		}

		_, err = s.DB.Exec(`
			DELETE FROM alias WHERE address='` + email + `'
		`)

		return err
	}

	isMailboxAlias, err := s.isMailboxAlias(email)
	if err != nil {
		return err
	}

	if isMailboxAlias {
		_, err = s.DB.Exec(`
			DELETE FROM forwardings WHERE address='` + email + `' AND is_alias=1
			`)

		return err
	}

	return fmt.Errorf("Alias %v does not exist", email)
}
