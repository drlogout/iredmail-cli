package iredmail

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	mailboxQueryByDomain   = "WHERE domain = ?"
	mailboxQueryAll        = ""
	mailboxQueryByUserName = "WHERE username = ?"
)

type MailboxEmail string

// types
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

type Mailboxes []Mailbox

func (m *Mailbox) IsCopyKept() bool {
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
		var username, password, name, domain, maildir string
		var quota int

		err := rows.Scan(&username, &password, &name, &domain, &quota, &maildir)
		if err != nil {
			return mailboxes, err
		}

		forwardings, err := s.forwardingsByMailbox(username)
		if err != nil {
			return mailboxes, err
		}

		mailboxAliases, err := s.queryMailboxAliases(username)
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
	return s.mailboxQuery(mailboxQueryAll)
}

func (s *Server) Mailbox(email string) (Mailbox, error) {
	mailbox := Mailbox{}

	exists, err := s.mailboxExists(email)
	if err != nil {
		return mailbox, err
	}

	if !exists {
		return mailbox, fmt.Errorf("Mailbox does not exist")
	}

	mailboxes, err := s.mailboxQuery(mailboxQueryByUserName, email)
	if err != nil {
		return mailbox, err
	}
	if len(mailboxes) == 0 {
		return mailbox, fmt.Errorf("Mailbox not found")
	}

	return mailboxes[0], nil
}

func (s *Server) MailboxAdd(email, password string, quota int, storageBasePath string) (Mailbox, error) {
	name, domain := parseEmail(email)
	m := Mailbox{
		Email:  email,
		Name:   name,
		Domain: domain,
		Quota:  quota,
	}

	domainExists, err := s.domainExists(domain)
	if err != nil {
		return m, err
	}
	if !domainExists {
		err := s.DomainAdd(Domain{
			Domain:   domain,
			Settings: DomainDefaultSettings,
		})
		if err != nil {
			return m, err
		}
	}

	mailboxExists, err := s.mailboxExists(email)
	if err != nil {
		return m, err
	}
	if mailboxExists {
		return m, fmt.Errorf("Mailbox %s already exists", email)
	}

	aliasExists, err := s.aliasExists(email)
	if err != nil {
		return m, err
	}
	if aliasExists {
		return m, fmt.Errorf("An alias %s already exists", email)
	}

	mailboxAliasExists, err := s.mailboxAliasExists(email)
	if err != nil {
		return m, err
	}
	if mailboxAliasExists {
		return m, fmt.Errorf("A mailbox alias %s already exists", email)
	}

	hash, err := generatePassword(password)
	if err != nil {
		return m, err
	}

	m.PasswordHash = hash

	mailDirHash := generateMaildirHash(email)
	storageBase := filepath.Dir(storageBasePath)
	storageNode := filepath.Base(storageBasePath)

	_, err = s.DB.Exec(`
		INSERT INTO mailbox (username, password, name,
			storagebasedirectory, storagenode, maildir,
			quota, domain, active, passwordlastchange, created)
		VALUES ('` + email + `', '` + hash + `', '` + name + `',
			'` + storageBase + `','` + storageNode + `', '` + mailDirHash + `',
			'` + strconv.Itoa(quota) + `', '` + domain + `', '1', NOW(), NOW());
		`)
	if err != nil {
		return m, err
	}

	err = s.ForwardingAdd(email, email)
	m.Forwardings = Forwardings{
		Forwarding{
			Address:    email,
			Forwarding: email,
		},
	}

	return m, err
}

func (s *Server) MailboxDelete(email string) error {
	mailboxExists, err := s.mailboxExists(email)
	if err != nil {
		return err
	}
	if !mailboxExists {
		return fmt.Errorf("Mailbox %s doesn't exist", email)
	}

	var mailDir string

	err = s.DB.QueryRow("SELECT maildir FROM mailbox WHERE username='" + email + "'").Scan(&mailDir)
	if err != nil {
		return err
	}

	err = os.RemoveAll(mailDir)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(`DELETE FROM mailbox WHERE username='` + email + `';`)
	if err != nil {
		return err
	}

	err = s.ForwardingDeleteAll(email)
	if err != nil {
		return err
	}

	err = s.MailboxAliasDeleteAll(email)

	return err
}

func (s *Server) MailboxUpdate(mailbox Mailbox) error {
	query := `
	UPDATE mailbox
	SET quota = ?, password = ?
	WHERE username = ?;`
	_, err := s.DB.Exec(query, mailbox.Quota, mailbox.PasswordHash, mailbox.Email)
	if err != nil {
		return err
	}

	return err
}

func (s *Server) MailboxKeepCopy(mailbox Mailbox, keepCopyInMailbox bool) error {
	if len(mailbox.Forwardings.External()) == 0 {
		return fmt.Errorf("No existing forwardings")
	}

	isCopyKept := mailbox.IsCopyKept()

	if isCopyKept && !keepCopyInMailbox {
		err := s.ForwardingDelete(mailbox.Email, mailbox.Email)
		if err != nil {
			return err
		}
	}

	if !isCopyKept && keepCopyInMailbox {
		err := s.ForwardingAdd(mailbox.Email, mailbox.Email)
		if err != nil {
			return err
		}
	}

	return nil
}
