package iredmail

import (
	"fmt"
)

func (s *Server) DomainAdd(domain Domain) error {
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
