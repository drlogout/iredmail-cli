package iredmail

import (
	"fmt"
	"strings"
)

type MailboxEmail string

// types
type Mailbox struct {
	Email          string
	Name           string
	Domain         string
	PasswordHash   string
	Quota          int
	Type           string
	MailDir        string
	MailboxAliases Forwardings
	Forwardings
}

type Mailboxes []Mailbox

func (m *Mailbox) IsCopyLeft() bool {
	for _, f := range m.Forwardings {
		if m.Email == f.Forwarding {
			return true
		}
	}
	return false
}

func (mailboxes Mailboxes) FilterBy(filter string) Mailboxes {
	filteredMailboxes := Mailboxes{}

	for _, m := range mailboxes {
		if strings.Contains(m.Email, filter) {
			filteredMailboxes = append(filteredMailboxes, m)
		}
	}

	return filteredMailboxes
}

func (s *Server) mailboxQuery(options queryOptions) (Mailboxes, error) {
	mailboxes := Mailboxes{}

	whereOption := ""
	if len(options.where) > 1 {
		whereOption = fmt.Sprintf("WHERE %v", options.where)
	}

	rows, err := s.DB.Query(`SELECT username, password, name, domain, quota, maildir FROM mailbox
` + whereOption + `
ORDER BY domain ASC, name ASC;`)
	if err != nil {
		return mailboxes, err
	}
	defer rows.Close()

	for rows.Next() {
		var username, password, name, domain, maildir string
		var quota int

		err := rows.Scan(&username, &password, &name, &domain, &quota, &maildir)
		if err != nil {
			return mailboxes, err
		}

		forwardings, err := s.queryForwardings(queryOptions{
			where: "address = '" + username + "' AND is_forwarding = 1",
		})
		if err != nil {
			return mailboxes, err
		}

		mailboxAliases, err := s.queryForwardings(queryOptions{
			where: "forwarding = '" + username + "' AND is_alias = 1",
		})
		if err != nil {
			return mailboxes, err
		}

		mailboxes = append(mailboxes, Mailbox{
			Email:          username,
			Name:           name,
			Domain:         domain,
			PasswordHash:   password,
			Quota:          quota,
			MailDir:        maildir,
			Forwardings:    forwardings,
			MailboxAliases: mailboxAliases,
		})
	}
	err = rows.Err()

	return mailboxes, err
}

func (s *Server) mailboxExists(email string) (bool, error) {
	var exists bool

	query := `SELECT exists
	(SELECT username FROM mailbox
		WHERE username = '` + email + `');`

	err := s.DB.QueryRow(query).Scan(&exists)
	if err != nil {
		return exists, err
	}

	if exists {
		return true, nil
	}

	return exists, nil
}

func (s *Server) Mailboxes() (Mailboxes, error) {
	return s.mailboxQuery(queryOptions{})
}

func (s *Server) Mailbox(email string) (Mailbox, error) {
	exists, err := s.mailboxExists(email)
	if err != nil {
		return Mailbox{}, err
	}

	if !exists {
		return Mailbox{}, fmt.Errorf("Mailbox does not exist")
	}

	mailboxes, err := s.mailboxQuery(queryOptions{
		where: `username = '` + email + `'`,
	})
	if err != nil {
		return Mailbox{}, err
	}
	if len(mailboxes) == 0 {
		return Mailbox{}, fmt.Errorf("Mailbox not found")
	}

	return mailboxes[0], nil
}
