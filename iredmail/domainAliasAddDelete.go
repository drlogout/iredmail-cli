package iredmail

import (
	"fmt"
)

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
