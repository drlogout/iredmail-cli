package iredmail

import (
	"fmt"
	"strings"
)

const (
	DomainDefaultSettings = "default_user_quota:2048"
)

type Domain struct {
	Domain      string
	Description string
	Settings    string
}

type Domains []Domain

func (d Domains) FilterBy(filter string) Domains {
	filteredDomains := Domains{}

	for _, domain := range d {
		if strings.Contains(domain.Domain, filter) {
			filteredDomains = append(filteredDomains, domain)
		}
	}

	return filteredDomains
}

func (s *Server) domainExists(domain string) (bool, error) {
	var exists bool

	sqlQuery := `
	SELECT exists
	(SELECT * FROM domain
	WHERE domain = ?);`

	err := s.DB.QueryRow(sqlQuery, domain).Scan(&exists)

	return exists, err
}

func (s *Server) Domains() (Domains, error) {
	domains := Domains{}

	rows, err := s.DB.Query(`SELECT domain, description, settings FROM domain ORDER BY domain ASC;`)
	if err != nil {
		return domains, err
	}
	defer rows.Close()

	for rows.Next() {
		var domain, description, settings string

		err := rows.Scan(&domain, &description, &settings)
		if err != nil {
			return domains, err
		}

		domains = append(domains, Domain{
			Domain:      domain,
			Description: description,
			Settings:    settings,
		})
	}
	err = rows.Err()

	return domains, err
}

func (s *Server) Domain(domainName string) (Domain, error) {
	var domain, description, settings string

	err := s.DB.QueryRow("SELECT domain, description, settings FROM domain WHERE domain =?", domainName).Scan(&domain, &description, &settings)
	if domain == "" {
		return Domain{}, fmt.Errorf("Domain %s doesn't exist", domainName)
	}

	d := Domain{
		Domain:      domain,
		Description: description,
		Settings:    settings,
	}

	return d, err
}

func (s *Server) DomainAdd(domain Domain) error {
	exists, err := s.domainExists(domain.Domain)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Domain %s already exists", domain.Domain)
	}

	sqlQuery := `
	INSERT INTO domain (domain, description, settings)
	VALUES (?, ?, ?);`

	_, err = s.DB.Exec(sqlQuery, domain.Domain, domain.Description, domain.Settings)

	return err
}

func (s *Server) DomainDelete(domain string) error {
	exists, err := s.domainExists(domain)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Domain %s doesn't exist", domain)
	}

	domainMailboxes, err := s.mailboxQuery(mailboxQueryByDomain, domain)
	if err != nil {
		return err
	}
	if len(domainMailboxes) > 0 {
		return fmt.Errorf("The domain %s still has mailboxes you need to delete them before", domain)
	}

	sqlQuery := `DELETE FROM domain WHERE domain = ?;`
	_, err = s.DB.Exec(sqlQuery, domain)

	return err
}
