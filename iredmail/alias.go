package iredmail

import (
	"fmt"
	"strings"
)

type Alias struct {
	Email  string
	Name   string
	Domain string
	Active bool
}

type Aliases []Alias

func (a Aliases) FilterBy(filter string) Aliases {
	filteredAliases := Aliases{}

	for _, al := range a {
		if strings.Contains(al.Email, filter) {
			filteredAliases = append(filteredAliases, al)
		}
	}

	return filteredAliases
}

func (s *Server) queryAliases(query string) (Aliases, error) {
	aliases := Aliases{}
	rows, err := s.DB.Query(query)
	if err != nil {
		return aliases, err
	}
	defer rows.Close()

	for rows.Next() {
		var address, domain string
		var active bool

		err := rows.Scan(&address, &domain, &active)
		if err != nil {
			return aliases, err
		}

		name, _ := parseEmail(address)

		aliases = append(aliases, Alias{
			Email:  address,
			Name:   name,
			Domain: domain,
			Active: active,
		})
	}
	err = rows.Err()

	return aliases, err
}

func (s *Server) AliasList() (Aliases, error) {
	return s.queryAliases(`SELECT address, domain, active FROM alias ORDER BY domain ASC, address ASC;`)
}

func (s *Server) AliasAdd(email string) error {
	_, domain := parseEmail(email)

	domainExists, err := s.DomainExists(domain)
	if err != nil {
		return err
	}
	if !domainExists {
		return fmt.Errorf("Domain %v does not exist", domain)
	}

	mailboxExists, err := s.MailboxExists(email)
	if err != nil {
		return err
	}
	if !mailboxExists {
		return fmt.Errorf("A mailbox with %v already exists", email)
	}

	_, err = s.DB.Exec(`
		REPLACE INTO alias (address, domain, active)
		VALUES ('` + email + `', '` + domain + `', 1)
	`)

	return err
}
