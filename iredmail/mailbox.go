package iredmail

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	mailboxQueryByDomain   = "WHERE domain = ?"
	mailboxQueryAll        = ""
	mailboxQueryByUserName = "WHERE username = ?"
)

// Mailbox struct
type Mailbox struct {
	Email          string
	Name           string
	Domain         string
	PasswordHash   string
	Quota          int
	MailDir        string
	MailboxAliases Forwardings
	Forwardings    Forwardings
}

// Mailboxes ...
type Mailboxes []Mailbox

// FilterBy is method that filters Mailboxes by a given string
func (mailboxes Mailboxes) FilterBy(filter string) Mailboxes {
	filteredMailboxes := Mailboxes{}

	for _, m := range mailboxes {
		if strings.Contains(m.Email, filter) {
			filteredMailboxes = append(filteredMailboxes, m)
		}
	}

	return filteredMailboxes
}

func (s *Server) mailboxQuery(whereQuery string, args ...interface{}) (Mailboxes, error) {
	mailboxes := Mailboxes{}

	sqlQuery := `SELECT username, password, name, domain, quota, maildir FROM mailbox
	` + whereQuery + `
	ORDER BY domain ASC, name ASC;`

	rows, err := s.DB.Query(sqlQuery, args...)
	if err != nil {
		return mailboxes, err
	}
	defer rows.Close()

	for rows.Next() {
		var mailboxEmail, password, name, domain, maildir string
		var quota int

		err := rows.Scan(&mailboxEmail, &password, &name, &domain, &quota, &maildir)
		if err != nil {
			return mailboxes, err
		}

		forwardings, err := s.forwardingsByMailbox(mailboxEmail)
		if err != nil {
			return mailboxes, err
		}

		mailboxAliases, err := s.mailboxAliaseQuery(mailboxEmail)
		if err != nil {
			return mailboxes, err
		}

		mailboxes = append(mailboxes, Mailbox{
			Email:          mailboxEmail,
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

func (s *Server) mailboxExists(mailboxEmail string) (bool, error) {
	var exists bool

	sqlQuery := `SELECT exists
	(SELECT username FROM mailbox
	WHERE username = ?);`
	err := s.DB.QueryRow(sqlQuery, mailboxEmail).Scan(&exists)

	return exists, err
}

// Mailboxes returns all mailboxes
func (s *Server) Mailboxes() (Mailboxes, error) {
	return s.mailboxQuery(mailboxQueryAll)
}

// Mailbox returns a Mailbox by mailboxEmail
func (s *Server) Mailbox(mailboxEmail string) (Mailbox, error) {
	mailbox := Mailbox{}

	exists, err := s.mailboxExists(mailboxEmail)
	if err != nil {
		return mailbox, err
	}
	if !exists {
		return mailbox, fmt.Errorf("Mailbox doesn't exist")
	}

	mailboxes, err := s.mailboxQuery(mailboxQueryByUserName, mailboxEmail)
	if err != nil {
		return mailbox, err
	}
	if len(mailboxes) == 0 {
		return mailbox, fmt.Errorf("Mailbox not found")
	}

	return mailboxes[0], nil
}

// MailboxAdd adds a new mailbox
func (s *Server) MailboxAdd(mailboxEmail, password string, quota int, storageBasePath string) error {
	name, domain := parseEmail(mailboxEmail)

	m := Mailbox{
		Email:  mailboxEmail,
		Name:   name,
		Domain: domain,
		Quota:  quota,
	}

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

	mailboxExists, err := s.mailboxExists(mailboxEmail)
	if err != nil {
		return err
	}
	if mailboxExists {
		return fmt.Errorf("Mailbox %s already exists", mailboxEmail)
	}

	aliasExists, err := s.aliasExists(mailboxEmail)
	if err != nil {
		return err
	}
	if aliasExists {
		return fmt.Errorf("An alias %s already exists", mailboxEmail)
	}

	mailboxAliasExists, err := s.mailboxAliasExists(mailboxEmail)
	if err != nil {
		return err
	}
	if mailboxAliasExists {
		return fmt.Errorf("A mailbox alias %s already exists", mailboxEmail)
	}

	hash, err := generatePassword(password)
	if err != nil {
		return err
	}

	m.PasswordHash = hash

	mailDirHash := generateMaildirHash(mailboxEmail)
	storageBase := filepath.Dir(storageBasePath)
	storageNode := filepath.Base(storageBasePath)

	sqlQuery := `
	INSERT INTO mailbox (username, password, name, storagebasedirectory, storagenode, maildir, quota, domain, active, passwordlastchange, created)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, '1', NOW(), NOW());`
	_, err = s.DB.Exec(sqlQuery, mailboxEmail, hash, name, storageBase, storageNode, mailDirHash, quota, domain)
	if err != nil {
		return err
	}

	err = s.ForwardingAdd(mailboxEmail, mailboxEmail)

	return err
}

// MailboxDelete delets a mailbox
func (s *Server) MailboxDelete(mailboxEmail string) error {
	mailboxExists, err := s.mailboxExists(mailboxEmail)
	if err != nil {
		return err
	}
	if !mailboxExists {
		return fmt.Errorf("Mailbox %s doesn't exist", mailboxEmail)
	}

	var mailDir string

	sqlQuery := "SELECT maildir FROM mailbox WHERE username = ?;"
	err = s.DB.QueryRow(sqlQuery, mailboxEmail).Scan(&mailDir)
	if err != nil {
		return err
	}

	err = os.RemoveAll(mailDir)
	if err != nil {
		return err
	}

	sqlQuery = "DELETE FROM mailbox WHERE username = ?;"
	_, err = s.DB.Exec(sqlQuery, mailboxEmail)
	if err != nil {
		return err
	}

	err = s.forwardingDeleteAll(mailboxEmail)
	if err != nil {
		return err
	}

	err = s.MailboxAliasDeleteAll(mailboxEmail)

	return err
}

// MailboxSetQuota sets the mailbox quota
func (s *Server) MailboxSetQuota(mailboxEmail string, quota int) error {
	sqlQuery := `UPDATE mailbox
	SET quota = ?
	WHERE username = ?;`
	_, err := s.DB.Exec(sqlQuery, quota, mailboxEmail)
	if err != nil {
		return err
	}

	return err
}

// MailboxSetKeepCopy sets the keep-copy behavior if forwardings exist
func (s *Server) MailboxSetKeepCopy(mailboxEmail string, keepCopyInMailbox bool) error {
	mailboxExists, err := s.mailboxExists(mailboxEmail)
	if err != nil {
		return err
	}
	if !mailboxExists {
		return fmt.Errorf("Mailbox %s doesn't exist", mailboxEmail)
	}

	mailbox, err := s.Mailbox(mailboxEmail)
	if err != nil {
		return err
	}

	if len(mailbox.Forwardings) == 0 {
		return fmt.Errorf("No forwardings exist for mailbox %s", mailboxEmail)
	}

	exists, err := s.forwardingExists(mailboxEmail, mailboxEmail)
	if err != nil {
		return err
	}

	if !keepCopyInMailbox {
		if !exists {
			return fmt.Errorf("keep-copy is already disabled")
		}

		err := s.ForwardingDelete(mailbox.Email, mailbox.Email)
		if err != nil {
			return err
		}
	} else {
		if exists {
			return fmt.Errorf("keep-copy is already enabled")
		}
		err := s.ForwardingAdd(mailbox.Email, mailbox.Email)
		if err != nil {
			return err
		}
	}

	return nil
}
