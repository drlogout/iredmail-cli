package iredmail

import (
	"fmt"
)

const (
	domainAliasQueryAll      = ""
	domainAliasQueryByDomain = "WHERE target_domain = ?"
)

// DomainAlias struct
type DomainAlias struct {
	Domain      string
	AliasDomain string
}

// DomainAliases ...
type DomainAliases []DomainAlias

func (s *Server) domainAliasQuery(whereQuery string, args ...interface{}) (DomainAliases, error) {
	aliasDomains := DomainAliases{}

	sqlQuery := `
	SELECT alias_domain, target_domain FROM alias_domain 
	` + whereQuery + `
	ORDER BY target_domain ASC;`

	rows, err := s.DB.Query(sqlQuery, args...)
	if err != nil {
		return aliasDomains, err
	}
	defer rows.Close()

	for rows.Next() {
		var aliasDomain, targetDomain string

		err := rows.Scan(&aliasDomain, &targetDomain)
		if err != nil {
			return aliasDomains, err
		}

		aliasDomains = append(aliasDomains, DomainAlias{
			AliasDomain: aliasDomain,
			Domain:      targetDomain,
		})
	}
	err = rows.Err()

	return aliasDomains, err
}

func (s *Server) domainAliasExists(aliasDomain string) (bool, error) {
	var exists bool

	query := `SELECT exists
	(SELECT * FROM alias_domain
	WHERE alias_domain = '` + aliasDomain + `');`

	err := s.DB.QueryRow(query).Scan(&exists)
	if err != nil {
		return exists, err
	}

	if exists {
		return true, nil
	}

	return exists, nil
}

// DomainAliasAdd adds a new domain alias
func (s *Server) DomainAliasAdd(aliasDomain, domain string) error {
	exists, err := s.domainExists(domain)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Domain %s doesn't exists", domain)
	}

	aliasExists, err := s.domainAliasExists(aliasDomain)
	if err != nil {
		return err
	}
	if aliasExists {
		return fmt.Errorf("Alias domain %s %s %s alreday exists", aliasDomain, arrowRight, domain)
	}

	_, err = s.DB.Exec(`
		INSERT INTO alias_domain (alias_domain, target_domain)
		VALUES ('` + aliasDomain + `', '` + domain + `')
	`)

	return err
}

// DomainAliasDelete deletes a domain alias
func (s *Server) DomainAliasDelete(aliasDomain string) error {
	aliasExists, err := s.domainAliasExists(aliasDomain)
	if err != nil {
		return err
	}
	if !aliasExists {
		return fmt.Errorf("Alias domain %s doesn't exist", aliasDomain)
	}

	_, err = s.DB.Exec(`DELETE FROM alias_domain WHERE alias_domain = '` + aliasDomain + `';`)

	return err
}

// DomainAliases returns all domainaliases
func (s *Server) DomainAliases() (DomainAliases, error) {
	return s.domainAliasQuery(domainAliasQueryAll)
}
