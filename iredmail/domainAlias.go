package iredmail

import (
	"fmt"
)

type AliasDomain struct {
	Domain      string
	AliasDomain string
}

type AliasDomains []AliasDomain

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
		return fmt.Errorf("Domain %v doesn't exists", targetDomain)
	}

	aliasExists, err := s.domainAliasExists(aliasDomain)
	if err != nil {
		return err
	}
	if aliasExists {
		return fmt.Errorf("Alias domain %v alreday exists", aliasDomain)
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
		return fmt.Errorf("Alias domain %v doesn't exist", aliasDomain)
	}

	_, err = s.DB.Exec(`DELETE FROM alias_domain WHERE alias_domain = '` + aliasDomain + `';`)

	return err
}

func (s *Server) DomainAliasList() (AliasDomains, error) {
	aliasDomains := AliasDomains{}

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

		aliasDomains = append(aliasDomains, AliasDomain{
			Domain:      targetDomain,
			AliasDomain: aliasDomain,
		})
	}
	err = rows.Err()

	return aliasDomains, err
}
