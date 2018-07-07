package iredmail

import (
	"fmt"
)

func (s *Server) domainCatchallExists(domain, catchallEmail string) (bool, error) {
	var exists bool

	query := `SELECT exists
	(SELECT forwarding FROM forwardings
	WHERE address = ? AND forwarding = ? AND is_forwarding = 0 AND is_alias = 0 AND is_list = 0);`

	err := s.DB.QueryRow(query, domain, catchallEmail).Scan(&exists)
	return exists, err
}

// DomainCatchallAdd adds a new catchall mailbox
func (s *Server) DomainCatchallAdd(domain, catchallEmail string) error {
	exists, err := s.domainExists(domain)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Domain %s doesn't exists", domain)
	}

	catchallExists, err := s.domainCatchallExists(domain, catchallEmail)
	if err != nil {
		return err
	}
	if catchallExists {
		return fmt.Errorf("Catch-all forwarding %s %s %s already exists", domain, arrowRight, catchallEmail)
	}

	_, forwardingDomain := parseEmail(catchallEmail)

	sqlQuery := `INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_forwarding, is_alias, is_list, active)
	VALUES (?, ?, ?, ?, 0, 0, 0, 1);`
	_, err = s.DB.Exec(sqlQuery, domain, catchallEmail, domain, forwardingDomain)

	return err
}

// DomainCatchallDelete deletes a domain alias
func (s *Server) DomainCatchallDelete(domain, catchallEmail string) error {
	domainExists, err := s.domainExists(domain)
	if err != nil {
		return err
	}
	if !domainExists {
		return fmt.Errorf("Domain %s doesn't exists", domain)
	}

	catchallExists, err := s.domainCatchallExists(domain, catchallEmail)
	if err != nil {
		return err
	}
	if !catchallExists {
		return fmt.Errorf("Catch-all forwarding %s %s %s doesn't exist", domain, arrowRight, catchallEmail)
	}

	sqlQuery := `DELETE FROM forwardings
	WHERE address = ? AND forwarding = ? AND is_forwarding = 0 AND is_alias = 0 AND is_list = 0;`
	_, err = s.DB.Exec(sqlQuery, domain, catchallEmail)

	return err
}

// domainCatchallDeleteAll deletes all catch-all forwardings of a domain
func (s *Server) domainCatchallDeleteAll(domain string) error {
	domainExists, err := s.domainExists(domain)
	if err != nil {
		return err
	}
	if !domainExists {
		return fmt.Errorf("Domain %s doesn't exists", domain)
	}

	sqlQuery := `DELETE FROM forwardings WHERE domain = ? AND is_forwarding = 0 AND is_alias = 0 AND is_list = 0;`
	_, err = s.DB.Exec(sqlQuery, domain)

	return err
}
