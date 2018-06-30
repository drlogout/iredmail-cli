package iredmail

import "fmt"

func (s *Server) AliasAdd(aliasEmail string) error {
	_, domain := parseEmail(aliasEmail)

	domainExists, err := s.domainExists(domain)
	if err != nil {
		return err
	}
	if !domainExists {
		err := s.DomainAdd(Domain{
			Domain:   domain,
			Settings: DomainDefaultSettings,
		})
		if err != nil {
			return err
		}
	}

	mailboxExists, err := s.mailboxExists(aliasEmail)
	if err != nil {
		return err
	}
	if mailboxExists {
		return fmt.Errorf("There is already a mailbox %v", aliasEmail)
	}

	isMailboxAlias, err := s.mailboxAliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if isMailboxAlias {
		return fmt.Errorf("There is already a mailbox alias %v ", aliasEmail)
	}

	isAlias, err := s.regularAliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if isAlias {
		return fmt.Errorf("There is already an alias %v", aliasEmail)
	}

	_, err = s.DB.Exec(`
		INSERT INTO alias (address, domain, active)
		VALUES ('` + aliasEmail + `', '` + domain + `', 1)
	`)

	return err
}
