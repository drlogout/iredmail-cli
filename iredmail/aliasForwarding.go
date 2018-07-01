package iredmail

import "fmt"

func (s *Server) AliasForwardingAdd(aliasEmail, forwardingEmail string) error {
	aliasExists, err := s.aliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if !aliasExists {
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

func (s *Server) AliasForwardingDelete(aliasEmail, forwardingEmail string) error {
	aliasExists, err := s.aliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if !aliasExists {
		return fmt.Errorf("An alias with %v doesn't exists", aliasEmail)
	}

	sqlQuery := `
	DELETE FROM forwardings WHERE address = ? AND forwarding = ? AND is_list = 1;`
	_, err = s.DB.Exec(sqlQuery, aliasEmail, forwardingEmail)

	return err
}
