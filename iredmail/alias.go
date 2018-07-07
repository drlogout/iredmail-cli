package iredmail

import (
	"fmt"
	"strings"
)

const (
	aliasQueryAll       = ""
	aliasQueryByAddress = "WHERE address = ?"
	aliasQueryByDomain  = "WHERE domain = ?"
)

// Alias struct
type Alias struct {
	Address string
	Domain  string
	Active  bool
	Forwardings
}

// Aliases ...
type Aliases []Alias

// FilterBy is method that filters Aliases by a given string
func (aliases Aliases) FilterBy(filter string) Aliases {
	filteredAliases := Aliases{}

	for _, a := range aliases {
		if strings.Contains(a.Address, filter) ||
			len(a.Forwardings.FilterBy(filter)) > 0 {
			filteredAliases = append(filteredAliases, a)
		}
	}

	return filteredAliases
}

func (s *Server) aliasQuery(whereQuery string, args ...interface{}) (Aliases, error) {
	aliases := Aliases{}

	sqlQuery := `
	SELECT address, domain, active FROM alias
	` + whereQuery + `
	ORDER BY domain ASC, address ASC;`
	rows, err := s.DB.Query(sqlQuery, args...)
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

func (s *Server) aliasExists(aliasEmail string) (bool, error) {
	var exists bool

	sqlQuery := `
	SELECT exists
	(SELECT address FROM alias
	WHERE address = ?);`
	err := s.DB.QueryRow(sqlQuery, aliasEmail).Scan(&exists)

	return exists, err
}

// Aliases returns all Aliases with its forwardings
func (s *Server) Aliases() (Aliases, error) {
	aliases, err := s.aliasQuery(aliasQueryAll)
	if err != nil {
		return aliases, err
	}

	allAliasForwardings, err := s.forwardingQuery(forwardingQueryAliasForwardingsAll)
	if err != nil {
		return aliases, err
	}

	for i, a := range aliases {
		for _, f := range allAliasForwardings {
			if f.Address == a.Address {
				aliases[i].Forwardings = append(aliases[i].Forwardings, f)
			}
		}
	}

	return aliases, nil
}

// Alias returns an Alias with its forwardings
func (s *Server) Alias(aliasEmail string) (Alias, error) {
	alias := Alias{}

	aliasExists, err := s.aliasExists(aliasEmail)
	if err != nil {
		return alias, err
	}
	if !aliasExists {
		return alias, fmt.Errorf("Alias %s doesn't exist", aliasEmail)
	}

	aliases, err := s.aliasQuery(aliasQueryByAddress, aliasEmail)
	if err != nil {
		return alias, err
	}

	if len(aliases) == 0 {
		return alias, fmt.Errorf("Alias not found")
	}

	alias = aliases[0]

	forwardings, err := s.forwardingQuery(forwardingQueryAliasForwardingsByAliasEmail, aliasEmail)
	if err != nil {
		return alias, err
	}

	alias.Forwardings = forwardings

	return alias, nil
}

// AliasAdd adds a new alias
func (s *Server) AliasAdd(aliasEmail string) error {
	mailboxExists, err := s.mailboxExists(aliasEmail)
	if err != nil {
		return err
	}
	if mailboxExists {
		return fmt.Errorf("There is already a mailbox %s", aliasEmail)
	}

	mailboxAliasExists, err := s.mailboxAliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if mailboxAliasExists {
		return fmt.Errorf("There is already a mailbox alias %s ", aliasEmail)
	}

	aliasExists, err := s.aliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if aliasExists {
		return fmt.Errorf("Alias %s already exists", aliasEmail)
	}

	_, domain := parseEmail(aliasEmail)

	domainExists, err := s.domainExists(domain)
	if err != nil {
		return err
	}
	if !domainExists {
		err := s.DomainAdd(Domain{
			Domain: domain,
		})
		if err != nil {
			return err
		}
	}

	sqlQuery := `
	INSERT INTO alias (address, domain, active)
	VALUES (?, ?, 1);`
	_, err = s.DB.Exec(sqlQuery, aliasEmail, domain)

	return err
}

// AliasDelete deletes an alias an its forwardings
func (s *Server) AliasDelete(aliasEmail string) error {
	aliasExists, err := s.aliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if !aliasExists {
		return fmt.Errorf("Alias %s does not exist", aliasEmail)
	}

	tx, err := s.DB.Begin()
	stmt1, err := tx.Prepare("DELETE FROM forwardings WHERE address='" + aliasEmail + "' AND is_forwarding = 0 AND is_alias = 0 AND is_list=1")
	_, err = stmt1.Exec()

	stmt2, err := tx.Prepare("DELETE FROM alias WHERE address='" + aliasEmail + "'")
	_, err = stmt2.Exec()

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()

	return err
}
