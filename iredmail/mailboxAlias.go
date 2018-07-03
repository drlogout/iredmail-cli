package iredmail

import (
	"fmt"
)

func (s *Server) queryMailboxAliases(mailboxEmail string) (Forwardings, error) {
	return s.forwardingQuery(forwardingQueryMailboxAliasForwardingsByAddress, mailboxEmail)
}

func (s *Server) MailboxAliasAdd(alias, email string) error {
	_, domain := parseEmail(email)
	a := fmt.Sprintf("%s@%s", alias, domain)

	mailboxExists, err := s.mailboxExists(a)
	if err != nil {
		return err
	}
	if mailboxExists {
		return fmt.Errorf("An mailbox with %s already exists", a)
	}

	aliasExists, err := s.aliasExists(a)
	if err != nil {
		return err
	}
	if aliasExists {
		return fmt.Errorf("An alias with %s already exists", a)
	}

	_, err = s.DB.Exec(`
		INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_alias, active)
		VALUES ('` + a + `', '` + email + `', '` + domain + `', '` + domain + `', 1, 1)
	`)

	return err
}

func (s *Server) MailboxAliasDelete(aliasEmail string) error {
	aliasExists, err := s.mailboxAliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if !aliasExists {
		return fmt.Errorf("An alias with %s doesn't exists", aliasEmail)
	}

	_, err = s.DB.Exec(`
		DELETE FROM forwardings WHERE address = '` + aliasEmail + `' AND is_alias = 1
	`)

	return err
}

func (s *Server) MailboxAliasDeleteAll(mailboxEmail string) error {
	_, err := s.DB.Exec(`
		DELETE FROM forwardings WHERE forwarding = '` + mailboxEmail + `' AND is_alias = 1
	`)

	return err
}

func (s *Server) mailboxAliasExists(aliasEmail string) (bool, error) {
	var exists bool

	sqlQuery := `
	SELECT exists
	(SELECT address FROM forwardings
	WHERE address = ? AND is_alias = 1);`

	err := s.DB.QueryRow(sqlQuery, aliasEmail).Scan(&exists)

	return exists, err
}
