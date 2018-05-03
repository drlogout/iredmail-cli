package iredmail

import (
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

func (s *Server) aliasExists(email string) (bool, error) {
	var exists bool

	isAlias, err := s.isAlias(email)
	if err != nil {
		return exists, err
	}

	isMailboxAlias, err := s.isMailboxAlias(email)
	if err != nil {
		return exists, err
	}

	return (isAlias || isMailboxAlias), nil
}

func (s *Server) isAlias(email string) (bool, error) {
	var exists bool

	query := `SELECT exists
	(SELECT address FROM alias
	WHERE address = '` + email + `');`

	err := s.DB.QueryRow(query).Scan(&exists)
	if err != nil {
		return exists, err
	}
	if exists {
		return true, nil
	}

	return exists, nil
}

func (s *Server) isMailboxAlias(email string) (bool, error) {
	var exists bool

	query := `SELECT exists
	(SELECT address FROM forwardings
	WHERE address = '` + email + `'
	AND is_alias=1);`

	err := s.DB.QueryRow(query).Scan(&exists)
	if err != nil {
		return exists, err
	}

	if exists {
		return true, nil
	}

	return exists, nil
}
