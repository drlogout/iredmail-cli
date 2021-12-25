package iredmail

import (
	"fmt"
)

func (s *Server) mailboxAliaseQuery(mailboxEmail string) (Forwardings, error) {
	return s.forwardingQuery(forwardingQueryMailboxAliasForwardingsByMailboxEamil, mailboxEmail)
}

func (s *Server) mailboxAliasExists(aliasEmail string) (bool, error) {
	var exists bool

	sqlQuery := `
	SELECT exists
	(SELECT address FROM forwardings
	WHERE address = ? AND is_forwarding = 0 AND is_alias = 1 AND is_list = 0);`
	err := s.DB.QueryRow(sqlQuery, aliasEmail).Scan(&exists)

	return exists, err
}

// MailboxAliasAdd adds a new mailbox alias
func (s *Server) MailboxAliasAdd(alias, mailboxEmail string) error {
	mailboxExists, err := s.mailboxExists(mailboxEmail)
	if err != nil {
		return err
	}
	if !mailboxExists {
		return fmt.Errorf("Mailbox %s doesn't exist", mailboxEmail)
	}

	_, domain := parseEmail(mailboxEmail)
	aliasEmail := alias
	_, aliasDomain := parseEmail(aliasEmail)

	mailboxAliasExists, err := s.mailboxAliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if mailboxAliasExists {
		return fmt.Errorf("A mailbox alias with %s already exists", aliasEmail)
	}

	mailboxWithSameEmailExists, err := s.mailboxExists(aliasEmail)
	if err != nil {
		return err
	}
	if mailboxWithSameEmailExists {
		return fmt.Errorf("A mailbox with %s already exists", aliasEmail)
	}

	aliasExists, err := s.aliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if aliasExists {
		return fmt.Errorf("An alias with %s already exists", aliasEmail)
	}

	sqlQuery := `
	INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_forwarding, is_alias, is_list, active)
	VALUES (?, ?, ?, ?, 0, 1, 0, 1)`
	_, err = s.DB.Exec(sqlQuery, aliasEmail, mailboxEmail, aliasDomain, domain)

	return err
}

// MailboxAliasDelete deletes a mailbox alias
func (s *Server) MailboxAliasDelete(aliasEmail string) error {
	aliasExists, err := s.mailboxAliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if !aliasExists {
		return fmt.Errorf("An alias with %s doesn't exists", aliasEmail)
	}

	sqlQuery := `DELETE FROM forwardings 
	WHERE address = ? AND is_forwarding = 0 AND is_alias = 1 AND is_list = 0;`
	_, err = s.DB.Exec(sqlQuery, aliasEmail)

	return err
}

// MailboxAliasDeleteAll deletes all mailbox aliases of a mailbox
func (s *Server) MailboxAliasDeleteAll(mailboxEmail string) error {
	sqlQuery := "DELETE FROM forwardings WHERE forwarding = ? AND is_forwarding = 0 AND is_alias = 1 AND is_list = 0;"
	_, err := s.DB.Exec(sqlQuery, mailboxEmail)

	return err
}
