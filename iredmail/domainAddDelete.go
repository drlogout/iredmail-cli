package iredmail

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

func (s *Server) DomainAdd(domain Domain) error {
	if !govalidator.IsDNSName(domain.Domain) {
		return fmt.Errorf("%v is no valid domain name", domain.Domain)
	}

	exists, err := s.domainExists(domain.Domain)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Domain %v already exists", domain)
	}

	_, err = s.DB.Exec(`
		INSERT INTO domain (domain, description, settings)
		VALUES ('` + domain.Domain + `', '` + domain.Description + `', '` + domain.Settings + `')
	`)

	return err
}

func (s *Server) DomainDelete(domain string, args ...bool) error {
	domainMailboxes, err := s.mailboxQuery(queryOptions{
		where: "domain = '" + domain + "'",
	})
	if err != nil {
		return err
	}
	if len(domainMailboxes) > 0 {
		return fmt.Errorf("The domain %v still has mailboxes you need to delete them before", domain)
	}

	_, err = s.DB.Exec(`DELETE FROM domain WHERE domain = '` + domain + `';`)

	return err
}
