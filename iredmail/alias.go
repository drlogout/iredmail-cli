package iredmail

import (
	"fmt"
	"strings"
)

type Alias struct {
	Address string
	Domain  string
	Active  bool
	Type    string
	Forwardings
}

type Aliases []Alias

func (a Aliases) FilterBy(filter string) Aliases {
	filteredAliases := Aliases{}

	for _, al := range a {
		if strings.Contains(al.Address, filter) {
			filteredAliases = append(filteredAliases, al)
		}
	}

	return filteredAliases
}

func (s *Server) queryAliases(options queryOptions) (Aliases, error) {
	aliases := Aliases{}

	whereOption := ""
	if len(options.where) > 1 {
		whereOption = fmt.Sprintf("WHERE %v", options.where)
	}

	rows, err := s.DB.Query(`SELECT address, domain, active FROM alias
` + whereOption + `
ORDER BY domain ASC, address ASC;`)
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

		aliases = append(aliases, Alias{
			Address: address,
			Domain:  domain,
			Active:  active,
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

	isMailboxAlias, err := s.isUserAlias(email)
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

func (s *Server) isUserAlias(email string) (bool, error) {
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
