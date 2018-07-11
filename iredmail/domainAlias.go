package iredmail

import (
	"fmt"
	"strings"
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

func (a DomainAliases) FilterBy(filter string) DomainAliases {
	filteredAliases := DomainAliases{}

	for _, alias := range a {
		if strings.Contains(alias.AliasDomain, filter) ||
			strings.Contains(alias.Domain, filter) {
			filteredAliases = append(filteredAliases, alias)
		}
	}

	return filteredAliases
}

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

	slqQuery := `SELECT exists
	(SELECT alias_domain FROM alias_domain
	WHERE alias_domain = ?);`
	err := s.DB.QueryRow(slqQuery, aliasDomain).Scan(&exists)

	return exists, err
}

// DomainAliasAdd adds a new domain alias
func (s *Server) DomainAliasAdd(aliasDomain, domain string) error {
	asDomainExists, err := s.domainExists(aliasDomain)
	if err != nil {
		return err
	}
	if asDomainExists {
		return fmt.Errorf("%s is already a domain", aliasDomain)
	}

	domainExists, err := s.domainExists(domain)
	if err != nil {
		return err
	}
	if !domainExists {
		return fmt.Errorf("Domain %s doesn't exists", domain)
	}

	aliasExists, err := s.domainAliasExists(aliasDomain)
	if err != nil {
		return err
	}
	if aliasExists {
		return fmt.Errorf("Alias domain %s %s %s already exists", aliasDomain, arrowRight, domain)
	}

	sqlQuery := `
	INSERT INTO alias_domain (alias_domain, target_domain)
	VALUES (?, ?);`
	_, err = s.DB.Exec(sqlQuery, aliasDomain, domain)

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

	sqlQuery := "DELETE FROM alias_domain WHERE alias_domain = ?;"
	_, err = s.DB.Exec(sqlQuery, aliasDomain)

	return err
}

func (s *Server) domainAliasDeleteAll(domain string) error {
	sqlQuery := "DELETE FROM alias_domain WHERE target_domain = ?;"
	_, err := s.DB.Exec(sqlQuery, domain)

	return err
}

// DomainAliases returns all domainaliases
func (s *Server) DomainAliases() (DomainAliases, error) {
	return s.domainAliasQuery(domainAliasQueryAll)
}
