package iredmail

import (
	"path/filepath"
	"strconv"
	"strings"
)

type Mailbox struct {
	Email        string
	Name         string
	Domain       string
	PasswordHash string
	Quota        int
}

type Mailboxes []Mailbox

func (m Mailboxes) FilterBy(filter string) Mailboxes {
	filteredMailboxes := Mailboxes{}

	for _, mailbox := range m {
		if strings.Contains(mailbox.Email, filter) {
			filteredMailboxes = append(filteredMailboxes, mailbox)
		}
	}

	return filteredMailboxes
}

func (s *Server) mailboxQuery(query string) (Mailboxes, error) {
	mailboxes := Mailboxes{}
	rows, err := s.DB.Query(query)
	if err != nil {
		return mailboxes, err
	}
	defer rows.Close()

	for rows.Next() {
		var username, password, name, domain string
		var quota int

		err := rows.Scan(&username, &password, &name, &domain, &quota)
		if err != nil {
			return mailboxes, err
		}

		mailboxes = append(mailboxes, Mailbox{
			Email:        username,
			Name:         name,
			Domain:       domain,
			PasswordHash: password,
			Quota:        quota,
		})
	}
	err = rows.Err()

	return mailboxes, err
}

func (s *Server) MailboxList() (Mailboxes, error) {
	mailboxes, err := s.mailboxQuery(`SELECT username, password, name, domain, quota FROM mailbox ORDER BY domain ASC, name ASC;`)

	return mailboxes, err
}

func (s *Server) MailboxAdd(email, password string, quota int, storageBasePath string) (Mailbox, error) {
	name, domain := parseEmail(email)
	m := Mailbox{
		Email:  email,
		Name:   name,
		Domain: domain,
		Quota:  quota,
	}

	domainExists, err := s.DomainExists(domain)
	if err != nil {
		return m, err
	}

	if !domainExists {
		err = s.DomainCreate(domain, quota)
		if err != nil {
			return m, err
		}
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
			storagebasedirectory,storagenode, maildir,
			quota, domain, active, passwordlastchange, created)
		VALUES ('` + email + `', '` + hash + `', '` + name + `',
			'` + storageBase + `','` + storageNode + `', '` + mailDirHash + `',
			'` + strconv.Itoa(quota) + `', '` + domain + `', '1', NOW(), NOW());
		`)

	err = s.ForwardingAdd(email, email)

	return m, err
}

func (s *Server) MailboxExists(email string) (bool, error) {
	var exists bool

	query := `select exists
	(select username from mailbox
	where username = '` + email + `');`

	err := s.DB.QueryRow(query).Scan(&exists)
	if err != nil {
		return exists, err
	}

	if exists {
		return true, nil
	}

	return exists, nil
}
