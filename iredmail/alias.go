package iredmail

import (
	"fmt"
	"strings"
)

type Alias struct {
	Address string
	Domain  string
	Active  bool
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

	rows, err := s.DB.Query(`
	SELECT address, domain, active FROM alias
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

	sqlQuery := `
	SELECT exists
	(SELECT address FROM alias
	WHERE address = ?);`

	err := s.DB.QueryRow(sqlQuery, email).Scan(&exists)

	return exists, err
}

func (s *Server) mailboxAliasExists(aliasEmail string) (bool, error) {
	var exists bool

	sqlQuery := `
	SELECT exists
	(SELECT address FROM forwardings
	WHERE address = ? AND is_alias = 1);`

	err := s.DB.QueryRow(sqlQuery, aliasEmail).Scan(&exists)

	return exists, err
}

func (s *Server) aliasForwardingExists(aliasEmail, forwardingEmail string) (bool, error) {
	var exists bool

	sqlQuery := `
	SELECT exists
	(SELECT address FROM forwardings
	WHERE address = ? AND forwarding = ? AND is_list = 1
	);`
	err := s.DB.QueryRow(sqlQuery, aliasEmail, forwardingEmail).Scan(&exists)

	return exists, err
}

func (s *Server) AliasDelete(aliasEmail string) error {
	aliasExists, err := s.aliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if !aliasExists {
		return fmt.Errorf("Alias %v does not exist", aliasEmail)
	}

	tx, err := s.DB.Begin()
	stmt1, err := tx.Prepare("DELETE FROM forwardings WHERE address='" + aliasEmail + "' and is_list=1")
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

func (s *Server) AliasAdd(aliasEmail string) error {
	_, domain := parseEmail(aliasEmail)

	domainExists, err := s.domainExists(domain)
	if err != nil {
		return err
	}
	if !domainExists {
		err := s.DomainAdd(Domain{
			Domain:   domain,
			Settings: DomainDefaultSettings,
		})
		if err != nil {
			return err
		}
	}

	mailboxExists, err := s.mailboxExists(aliasEmail)
	if err != nil {
		return err
	}
	if mailboxExists {
		return fmt.Errorf("There is already a mailbox %v", aliasEmail)
	}

	isMailboxAlias, err := s.mailboxAliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if isMailboxAlias {
		return fmt.Errorf("There is already a mailbox alias %v ", aliasEmail)
	}

	aliasExists, err := s.aliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if aliasExists {
		return fmt.Errorf("There is already an alias %v", aliasEmail)
	}

	_, err = s.DB.Exec(`
		INSERT INTO alias (address, domain, active)
		VALUES ('` + aliasEmail + `', '` + domain + `', 1)
	`)

	return err
}

func (s *Server) Aliases() (Aliases, error) {
	aliases, err := s.queryAliases(queryOptions{})
	if err != nil {
		return aliases, err
	}

	sqlQuery := `
	SELECT address, domain, forwarding, dest_domain, active, is_alias, is_forwarding, is_list 
	FROM forwardings
	WHERE is_list = 1
	ORDER BY domain ASC, address ASC;`

	aliasForwardings, err := s.queryForwardings(sqlQuery)
	if err != nil {
		return aliases, err
	}

	for i, a := range aliases {
		for _, f := range aliasForwardings {
			if f.Address == a.Address {
				aliases[i].Forwardings = append(aliases[i].Forwardings, f)
			}
		}
	}

	return aliases, nil
}

func (s *Server) Alias(aliasEmail string) (Alias, error) {
	alias := Alias{}

	aliasExists, err := s.aliasExists(aliasEmail)
	if err != nil {
		return alias, err
	}
	if !aliasExists {
		return alias, fmt.Errorf("Alias %v doesn't exist", aliasEmail)
	}

	aliases, err := s.queryAliases(queryOptions{
		where: "address = '" + aliasEmail + "'",
	})
	if err != nil {
		return alias, err
	}

	if len(aliases) == 0 {
		return alias, fmt.Errorf("Alias not found")
	}

	alias = aliases[0]

	sqlQuery := `
	SELECT address, domain, forwarding, dest_domain, active, is_alias, is_forwarding, is_list 
	FROM forwardings
	WHERE forwarding = ? AND is_alias = 1
	ORDER BY domain ASC, address ASC;`

	forwardings, err := s.queryForwardings(sqlQuery, mailboxEmail)
	if err != nil {
		return alias, err
	}

	alias.Forwardings = forwardings

	return alias, nil
}
