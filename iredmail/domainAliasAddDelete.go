package iredmail

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

func (s *Server) DomainAliasAdd(aliasDomain, domain string) error {
	if !govalidator.IsDNSName(aliasDomain) {
		return fmt.Errorf("%v is no valid domain name", aliasDomain)
	}
	if !govalidator.IsDNSName(domain) {
		return fmt.Errorf("%v is no valid domain name", domain)
	}

	exists, err := s.domainExists(domain)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Domain %v doesn't exists", domain)
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
		VALUES ('` + aliasDomain + `', '` + domain + `')
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
