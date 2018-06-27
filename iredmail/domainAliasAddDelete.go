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

	aliasExists, err := s.domainAliasExists(aliasDomain, domain)
	if err != nil {
		return err
	}
	if aliasExists {
		return fmt.Errorf("Alias domain %v -> %v alreday exists", aliasDomain, domain)
	}

	_, err = s.DB.Exec(`
		INSERT INTO alias_domain (alias_domain, target_domain)
		VALUES ('` + aliasDomain + `', '` + domain + `')
	`)

	return err
}

func (s *Server) DomainAliasDelete(domain string, args ...bool) error {
	domainUsers, err := s.userQuery(queryOptions{
		where: "domain = '" + domain + "'",
	})
	if err != nil {
		return err
	}
	if len(domainUsers) > 0 {
		return fmt.Errorf("The domain %v still has users you need to delete them before", domain)
	}

	_, err = s.DB.Exec(`DELETE FROM domain WHERE domain = '` + domain + `';`)

	return err
}
