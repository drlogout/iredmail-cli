package iredmail

import "fmt"

func (s *Server) AliasAdd(email, destEmail string) error {
	_, domain := parseEmail(email)
	_, destDomain := parseEmail(destEmail)

	domainExists, err := s.DomainExists(domain)
	if err != nil {
		return err
	}
	if !domainExists {
		return fmt.Errorf("Domain %v does not exist, please create one first", domain)
	}

	mailboxExists, err := s.MailboxExists(email)
	if err != nil {
		return err
	}
	if mailboxExists {
		return fmt.Errorf("A mailbox %v already exists", email)
	}

	isAlias, err := s.isAlias(email)
	if err != nil {
		return err
	}
	if isAlias {
		return fmt.Errorf("An alias %v already exists", email)
	}

	isMailboxAlias, err := s.isMailboxAlias(email)
	if err != nil {
		return err
	}
	if isMailboxAlias {
		return fmt.Errorf("An alias %v already exists", email)
	}

	// If domain equals destDomain and destEmail is a local mailbox create mailbox alias
	destEmailIsMailbox, err := s.MailboxExists(destEmail)
	if err != nil {
		return err
	}
	if domain == destDomain && destEmailIsMailbox {
		_, err = s.DB.Exec(`
			INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_alias, active)
			VALUES ('` + email + `', '` + destEmail + `', '` + domain + `', '` + destDomain + `', 1, 1)
		`)

		return err
	}

	if !isAlias {
		_, err = s.DB.Exec(`
			INSERT INTO alias (address, domain, active)
			VALUES ('` + email + `', '` + domain + `', 1)
		`)
		if err != nil {
			return err
		}
	}

	_, err = s.DB.Exec(`
		INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_list, active)
		VALUES ('` + email + `', '` + destEmail + `', '` + domain + `', '` + destDomain + `', 1, 1)
	`)

	return err
}
