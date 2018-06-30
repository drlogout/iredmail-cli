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

func (s *Server) AliasAddForwarding(aliasEmail, forwardingEmail string) error {
	regularAliasExists, err := s.regularAliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if !regularAliasExists {
		return fmt.Errorf("Alias %v doesn't exist", aliasEmail)
	}

	forwardingExists, err := s.aliasForwardingExists(aliasEmail, forwardingEmail)
	if err != nil {
		return err
	}
	if forwardingExists {
		return fmt.Errorf("Alias forwarding %v %v %v already exists", aliasEmail, arrowRight, forwardingEmail)
	}

	_, aliasDomain := parseEmail(aliasEmail)
	_, forwardingDomain := parseEmail(forwardingEmail)

	sqlQuery := `
	INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_list, active)
	VALUES (?, ?, ?, ?, 1, 1);`
	_, err = s.DB.Exec(sqlQuery, aliasEmail, forwardingEmail, aliasDomain, forwardingDomain)

	return err
}
