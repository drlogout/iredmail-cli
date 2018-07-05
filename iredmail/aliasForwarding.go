package iredmail

import "fmt"

func (s *Server) aliasForwardingExists(aliasEmail, forwardingEmail string) (bool, error) {
	var exists bool

	sqlQuery := `
	SELECT exists
	(SELECT address FROM forwardings
	WHERE address = ? AND forwarding = ? AND is_list = 1 AND is_forwarding = 0 AND is_alias = 0);`
	err := s.DB.QueryRow(sqlQuery, aliasEmail, forwardingEmail).Scan(&exists)

	return exists, err
}

// AliasForwardingAdd adds a new alias fowarding
func (s *Server) AliasForwardingAdd(aliasEmail, forwardingEmail string) error {
	aliasExists, err := s.aliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if !aliasExists {
		return fmt.Errorf("Alias %s doesn't exist", aliasEmail)
	}

	forwardingExists, err := s.aliasForwardingExists(aliasEmail, forwardingEmail)
	if err != nil {
		return err
	}
	if forwardingExists {
		return fmt.Errorf("Alias forwarding %s %s %s already exists", aliasEmail, arrowRight, forwardingEmail)
	}

	_, aliasDomain := parseEmail(aliasEmail)
	_, forwardingDomain := parseEmail(forwardingEmail)

	sqlQuery := `
	INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_forwarding, is_alias, is_list, active)
	VALUES (?, ?, ?, ?, 0, 0, 1, 1);`
	_, err = s.DB.Exec(sqlQuery, aliasEmail, forwardingEmail, aliasDomain, forwardingDomain)

	return err
}

// AliasForwardingDelete deletes an alias forwarding
func (s *Server) AliasForwardingDelete(aliasEmail, forwardingEmail string) error {
	aliasExists, err := s.aliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if !aliasExists {
		return fmt.Errorf("An alias with %s doesn't exists", aliasEmail)
	}

	forwardingExists, err := s.aliasForwardingExists(aliasEmail, forwardingEmail)
	if err != nil {
		return err
	}
	if !forwardingExists {
		return fmt.Errorf("An alias forwarding %s %s %s doesn't exists", aliasEmail, arrowRight, forwardingEmail)
	}

	sqlQuery := `
	DELETE FROM forwardings WHERE address = ? AND forwarding = ? AND is_list = 1 AND is_forwarding = 0 AND is_alias = 0;`
	_, err = s.DB.Exec(sqlQuery, aliasEmail, forwardingEmail)

	return err
}
