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

type AliasDomain struct {
	Domain      string
	AliasDomain string
}

type AliasDomains []AliasDomain

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
