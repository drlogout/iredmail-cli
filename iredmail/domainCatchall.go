package iredmail

import (
	"fmt"
	"strings"
)

type Catchalls Forwardings

func (c Catchalls) FilterBy(filter string) Catchalls {
	filteredCatchalls := Catchalls{}

	for _, catchall := range c {
		if strings.Contains(catchall.Address, filter) ||
			strings.Contains(catchall.Forwarding, filter) {
			filteredCatchalls = append(filteredCatchalls, catchall)
		}
	}

	return filteredCatchalls
}

func (s *Server) domainCatchallExists(catchallEmail string) (bool, error) {
	var exists bool

	query := `SELECT exists
	(SELECT forwarding FROM forwardings
	WHERE forwarding = ? AND is_forwarding = 0 AND is_alias = 0 AND is_list = 0);`

	err := s.DB.QueryRow(query, catchallEmail).Scan(&exists)
	if err != nil {
		return exists, err
	}

	if exists {
		return true, nil
	}

	return exists, nil
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

	catchallExists, err := s.domainCatchallExists(catchallEmail)
	if err != nil {
		return err
	}
	if catchallExists {
		return fmt.Errorf("Catch-all mailbox %s already exists", catchallEmail)
	}

	_, forwardingDomain := parseEmail(catchallEmail)

	sqlQuery := `INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_forwarding, is_alias, is_list, active)
	VALUES (?, ?, ?, ?, 0, 0, 0, 1);`
	_, err = s.DB.Exec(sqlQuery, domain, catchallEmail, domain, forwardingDomain)

	return err
}

// // DomainAliasDelete deletes a domain alias
// func (s *Server) DomainAliasDelete(aliasDomain string) error {
// 	aliasExists, err := s.domainAliasExists(aliasDomain)
// 	if err != nil {
// 		return err
// 	}
// 	if !aliasExists {
// 		return fmt.Errorf("Alias domain %s doesn't exist", aliasDomain)
// 	}

// 	_, err = s.DB.Exec(`DELETE FROM alias_domain WHERE alias_domain = '` + aliasDomain + `';`)

// 	return err
// }

// func (s *Server) domainAliasDeleteAll(domain string) error {
// 	sqlQuery := `DELETE FROM alias_domain WHERE target_domain = ?;`
// 	_, err := s.DB.Exec(sqlQuery, domain)

// 	return err
// }

// // DomainAliases returns all domainaliases
// func (s *Server) DomainAliases() (DomainAliases, error) {
// 	return s.domainAliasQuery(domainAliasQueryAll)
// }
