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

func (s *Server) DomainExists(domain string) (bool, error) {
	var exists bool

	query := `select exists
	(select domain from domain
	where domain = '` + domain + `');`

	err := s.DB.QueryRow(query).Scan(&exists)
	if err != nil {
		return exists, err
	}

	if exists {
		return true, nil
	}

	return exists, nil
}

func (s *Server) DomainAdd(domain Domain) error {
	exists, err := s.DomainExists(domain.Domain)
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

func (s *Server) DomainList() (Domains, error) {
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

func (s *Server) DomainGet(domainName string) (Domain, error) {
	var domain, description, settings string

	err := s.DB.QueryRow("SELECT domain, description, settings FROM domain WHERE domain =?", domainName).Scan(&domain, &description, &settings)
	if domain == "" {
		return Domain{}, fmt.Errorf("Domain %v doesn't exist", domainName)
	}

	d := Domain{
		Domain:      domain,
		Description: description,
		Settings:    settings,
	}

	return d, err
}

func (s *Server) DomainUpdate(d Domain) error {
	_, err := s.DomainGet(d.Domain)
	if err != nil {
		return err
	}

	query := "UPDATE domain\n"

	if d.Description != "" && d.Settings != "" {
		query = query + " SET description='" + d.Description + "', settings='" + d.Settings + "'\n"
	}

	if d.Description != "" {
		query = query + " SET description='" + d.Description + "'\n"
	}

	if d.Settings != "" {
		query = query + " SET settings='" + d.Settings + "'\n"
	}

	query = query + " WHERE domain='" + d.Domain + "';"

	fmt.Println(query)

	_, err = s.DB.Exec(query)

	return err
}

type DomainInfo struct {
	Domain  Domain
	Users   Users
	Aliases Aliases
}

func (s *Server) DomainInfo(domainName string) error {

	return nil
}
