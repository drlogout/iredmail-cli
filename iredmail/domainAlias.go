package iredmail

import (
	"fmt"
)

const (
	domainAliasQueryAll      = ""
	domainAliasQueryByDomain = "WHERE target_domain = ?"
)

type DomainAlias struct {
	Domain      string
	AliasDomain string
}

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

func (s *Server) DomainAliasAdd(aliasDomain, targetDomain string) error {
	exists, err := s.domainExists(targetDomain)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Domain %s doesn't exists", targetDomain)
	}

	aliasExists, err := s.domainAliasExists(aliasDomain)
	if err != nil {
		return err
	}
	if aliasExists {
		return fmt.Errorf("Alias domain %s alreday exists", aliasDomain)
	}

	_, err = s.DB.Exec(`
		INSERT INTO alias_domain (alias_domain, target_domain)
		VALUES ('` + aliasDomain + `', '` + targetDomain + `')
	`)

	return err
}

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

func (s *Server) DomainAliasList() (DomainAliases, error) {
	aliasDomains := DomainAliases{}

	rows, err := s.DB.Query(`SELECT alias_domain, target_domain FROM alias_domain ORDER BY target_domain ASC;`)
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
			Domain:      targetDomain,
			AliasDomain: aliasDomain,
		})
	}
	err = rows.Err()

	return aliasDomains, err
}
