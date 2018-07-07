package iredmail

import (
	"fmt"
	"strings"
)

const (
	forwardingQueryCatchallForwardingsAll     = "WHERE is_forwarding = 0 AND is_alias = 0 AND is_list = 0"
	forwardingQueryForwardingsAll             = "WHERE is_forwarding = 1 AND is_alias = 0 AND is_list = 0"
	forwardingQueryMailboxAliasForwardingsAll = "WHERE is_forwarding = 0 AND is_alias = 1 AND is_list = 0"
	forwardingQueryAliasForwardingsAll        = "WHERE is_forwarding = 0 AND is_alias = 0 AND is_list = 1"

	forwardingQueryCatchallByDomain                      = "WHERE address = ? AND is_forwarding = 0 AND is_alias = 0 AND is_list = 0"
	forwardingQueryForwardingsByMailboxEmail             = "WHERE address = ? AND is_forwarding = 1 AND is_alias = 0 AND is_list = 0"
	forwardingQueryMailboxAliasForwardingsByMailboxEamil = "WHERE forwarding = ? AND is_forwarding = 0 AND is_alias = 1 AND is_list = 0"
	forwardingQueryAliasForwardingsByAliasEmail          = "WHERE address = ? AND is_forwarding = 0 AND is_alias = 0 AND is_list = 1"
)

// Forwarding struct
type Forwarding struct {
	Address             string
	Domain              string
	Forwarding          string
	DestDomain          string
	Active              bool
	IsList              bool
	IsAlias             bool
	IsForwarding        bool
	IsCopyLeftInMailbox bool
}

// Forwardings ...
type Forwardings []Forwarding

// FilterBy is method that filters Forwardings by a given string
func (forwardings Forwardings) FilterBy(filter string) Forwardings {
	filteredForwardings := Forwardings{}

	for _, f := range forwardings {
		if strings.Contains(f.Address, filter) ||
			strings.Contains(f.Forwarding, filter) {
			filteredForwardings = append(filteredForwardings, f)
		}
	}

	return filteredForwardings
}

func (s *Server) forwardingQuery(whereQuery string, args ...interface{}) (Forwardings, error) {
	Forwardings := Forwardings{}

	sqlQuery := `SELECT address, domain, forwarding, dest_domain, is_forwarding, is_alias, is_list, active
	FROM forwardings
	` + whereQuery + `
	ORDER BY domain ASC, address ASC;`
	rows, err := s.DB.Query(sqlQuery, args...)
	if err != nil {
		return Forwardings, err
	}
	defer rows.Close()

	for rows.Next() {
		var mailboxEmail, domain, destinationEmail, destinationDomain string
		var isForwarding, isAlias, isList, active bool

		err := rows.Scan(&mailboxEmail, &domain, &destinationEmail, &destinationDomain, &isForwarding, &isAlias, &isList, &active)
		if err != nil {
			return Forwardings, err
		}

		Forwardings = append(Forwardings, Forwarding{
			Address:      mailboxEmail,
			Domain:       domain,
			Forwarding:   destinationEmail,
			DestDomain:   destinationDomain,
			Active:       active,
			IsAlias:      isAlias,
			IsForwarding: isForwarding,
			IsList:       isList,
		})
	}
	err = rows.Err()

	return Forwardings, err
}

func (s *Server) forwardingExists(mailboxEmail, destinationEmail string) (bool, error) {
	var exists bool

	sqlQuery := `
	SELECT exists
	(SELECT address FROM forwardings
	WHERE address = ? AND forwarding = ? AND is_forwarding = 1 AND is_alias = 0 AND is_list = 0);`
	err := s.DB.QueryRow(sqlQuery, mailboxEmail, destinationEmail).Scan(&exists)

	return exists, err
}

// Forwardings returns all forwardings (actual forwardings, without mailbox copy)
func (s *Server) Forwardings() (Forwardings, error) {
	withoutMailboxCopy := Forwardings{}

	forwardings, err := s.forwardingQuery(forwardingQueryForwardingsAll)
	if err != nil {
		return withoutMailboxCopy, err
	}

	copyLeftInMailbox := map[string]bool{}

	for _, f := range forwardings {
		if _, ok := copyLeftInMailbox[f.Address]; !ok &&
			f.Address == f.Forwarding {
			copyLeftInMailbox[f.Address] = true
		}
	}

	for _, f := range forwardings {
		f.IsCopyLeftInMailbox = copyLeftInMailbox[f.Address]
		if f.Address != f.Forwarding {
			withoutMailboxCopy = append(withoutMailboxCopy, f)
		}
	}

	return withoutMailboxCopy, err
}

// forwardingsByMailbox returns all forwardings of a mailbox
func (s *Server) forwardingsByMailbox(mailboxEmail string) (Forwardings, error) {
	withoutMailboxCopy := Forwardings{}

	forwardings, err := s.forwardingQuery(forwardingQueryForwardingsByMailboxEmail, mailboxEmail)
	if err != nil {
		return withoutMailboxCopy, err
	}

	copyLeftInMailbox := map[string]bool{}

	for _, f := range forwardings {
		if _, ok := copyLeftInMailbox[f.Address]; !ok &&
			f.Address == f.Forwarding {
			copyLeftInMailbox[f.Address] = true
		}
	}

	for _, f := range forwardings {
		f.IsCopyLeftInMailbox = copyLeftInMailbox[f.Address]
		if f.Address != f.Forwarding {
			withoutMailboxCopy = append(withoutMailboxCopy, f)
		}
	}

	return withoutMailboxCopy, err
}

// ForwardingAdd adds a new Forwarding
func (s *Server) ForwardingAdd(mailboxEmail, destinationEmail string) error {
	mailboxExists, err := s.mailboxExists(mailboxEmail)
	if err != nil {
		return err
	}
	if !mailboxExists {
		return fmt.Errorf("Mailbox %s doesn't exist", mailboxEmail)
	}

	forwardingExists, err := s.forwardingExists(mailboxEmail, destinationEmail)
	if err != nil {
		return err
	}
	if forwardingExists {
		return fmt.Errorf("Forwarding %s %s %s already exists", mailboxEmail, arrowRight, destinationEmail)
	}

	_, domain := parseEmail(mailboxEmail)
	_, destinationDomain := parseEmail(destinationEmail)

	sqlQuery := `
	INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_forwarding, is_alias, is_list, active)
	VALUES (?, ?, ?, ?, 1, 0, 0, 1);`
	_, err = s.DB.Exec(sqlQuery, mailboxEmail, destinationEmail, domain, destinationDomain)

	return err
}

// ForwardingDelete deletes a forwarding
func (s *Server) ForwardingDelete(mailboxEmail, destinationEmail string) error {
	exists, err := s.forwardingExists(mailboxEmail, destinationEmail)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Forwarding %s %s %s doesn't exist", mailboxEmail, arrowRight, destinationEmail)
	}

	sqlQuery := `DELETE FROM forwardings 
	WHERE address = ? AND forwarding = ? AND is_forwarding = 1 AND is_list = 0 AND is_alias = 0;`
	_, err = s.DB.Exec(sqlQuery, mailboxEmail, destinationEmail)

	return err
}

// ForwardingDeleteAll deletes all forwardings of a mailbox
func (s *Server) ForwardingDeleteAll(mailboxEmail string) error {
	sqlQuery := `DELETE FROM forwardings 
	WHERE address = ? AND is_forwarding = 1 AND is_list = 0 AND is_alias = 0;`
	_, err := s.DB.Exec(sqlQuery, mailboxEmail)

	return err
}
