package iredmail

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

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
		return m, fmt.Errorf("Mailbox %v already exists", email)
	}

	aliasExists, err := s.aliasExists(email)
	if err != nil {
		return m, err
	}
	if aliasExists {
		return m, fmt.Errorf("An alias %v already exists", email)
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
			Mailbox:    email,
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
		return fmt.Errorf("Mailbox %v doesn't exist", email)
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

func (s *Server) MailboxSet(mailbox Mailbox) error {
	query := `
	UPDATE mailbox
	SET quota = ?
	WHERE username = ?;`
	_, err := s.DB.Exec(query, mailbox.Quota, mailbox.Email)
	if err != nil {
		return err
	}

	return err
}
