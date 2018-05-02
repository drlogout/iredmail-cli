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

func (s *Server) AliasAdd(email, destEmail string) error {
	_, domain := parseEmail(email)
	_, destDomain := parseEmail(destEmail)

	domainExists, err := s.DomainExists(domain)
	if err != nil {
		return err
	}
	if !domainExists {
		return fmt.Errorf("Domain %v does not exist, please create one first", domain)
	}

	mailboxExists, err := s.MailboxExists(email)
	if err != nil {
		return err
	}
	if mailboxExists {
		return fmt.Errorf("A mailbox %v already exists", email)
	}

	aliasExists, err := s.aliasExists(email)
	if err != nil {
		return err
	}
	if aliasExists {
		return fmt.Errorf("An alias %v already exists", email)
	}

	mailboxAliasExists, err := s.mailboxAliasExists(email)
	if err != nil {
		return err
	}
	if mailboxAliasExists {
		return fmt.Errorf("An alias %v already exists", email)
	}

	// If domain equals detsDomain and destEmail is a ocal mailbox create mailbox alias
	destEmailIsMailbox, err := s.MailboxExists(destEmail)
	if err != nil {
		return err
	}
	if domain == destDomain && destEmailIsMailbox {
		_, err = s.DB.Exec(`
			INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_alias, active)
			VALUES ('` + email + `', '` + destEmail + `', '` + domain + `', '` + destDomain + `', 1, 1)
		`)

		return err
	}

	if !aliasExists {
		_, err = s.DB.Exec(`
			INSERT INTO alias (address, domain, active)
			VALUES ('` + email + `', '` + domain + `', 1)
		`)
		if err != nil {
			return err
		}
	}

	_, err = s.DB.Exec(`
		INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_list, active)
		VALUES ('` + email + `', '` + destEmail + `', '` + domain + `', '` + destDomain + `', 1, 1)
	`)

	return err
}

func (s *Server) aliasExists(email string) (bool, error) {
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

func (s *Server) mailboxAliasExists(email string) (bool, error) {
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

func (s *Server) AliasRemove(email string) error {
	isAlias, err := s.aliasExists(email)
	if err != nil {
		return err
	}

	if isAlias {
		_, err = s.DB.Exec(`
			DELETE FROM forwardings WHERE address='` + email + `' and is_list=1
		`)
		if err != nil {
			return err
		}

		_, err = s.DB.Exec(`
			DELETE FROM alias WHERE address='` + email + `'
		`)

		return err
	}

	isMailboxAlias, err := s.mailboxAliasExists(email)
	if err != nil {
		return err
	}

	if isMailboxAlias {
		_, err = s.DB.Exec(`
			DELETE FROM forwardings WHERE address='` + email + `' AND is_alias=1
			`)
		return err
	}

	return fmt.Errorf("Alias %v does not exist", email)
}

func (s *Server) aliasCheck() error {
	result := []string{}

	forwardings, err := s.queryForwardings(queryOptions{})
	if err != nil {
		return err
	}

	for _, forwarding := range forwardings {
		// destEmailIsMailbox
		aliasExists, err := s.aliasExists(forwarding.Address)
		if err != nil {
			return err
		}

		fmt.Println(forwarding.Address, forwarding.IsAlias)
		if forwarding.IsAlias && aliasExists {
			result = append(result, fmt.Sprintf("%v Should be per-user alias address", forwarding.Address))
		}
	}

	for _, r := range result {
		fmt.Println(r)
	}
	return nil
}
