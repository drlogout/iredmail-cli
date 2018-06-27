package iredmail

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func (s *Server) UserAdd(email, password string, quota int, storageBasePath string) (User, error) {
	name, domain := parseEmail(email)
	u := User{
		Email:  email,
		Name:   name,
		Domain: domain,
		Quota:  quota,
	}

	domainExists, err := s.domainExists(domain)
	if err != nil {
		return u, err
	}
	if !domainExists {
		err := s.DomainAdd(Domain{
			Domain:   domain,
			Settings: DomainDefaultSettings,
		})
		if err != nil {
			return u, err
		}
	}

	userExists, err := s.userExists(email)
	if err != nil {
		return u, err
	}
	if userExists {
		return u, fmt.Errorf("User %v already exists", email)
	}

	aliasExists, err := s.aliasExists(email)
	if err != nil {
		return u, err
	}
	if aliasExists {
		return u, fmt.Errorf("An alias %v already exists", email)
	}

	hash, err := generatePassword(password)
	if err != nil {
		return u, err
	}

	u.PasswordHash = hash

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
		return u, err
	}

	err = s.ForwardingAdd(email, email)
	u.Forwardings = Forwardings{
		Forwarding{
			Address:    email,
			Forwarding: email,
		},
	}

	return u, err
}

func (s *Server) UserDelete(email string) error {
	userExists, err := s.userExists(email)
	if err != nil {
		return err
	}
	if !userExists {
		return fmt.Errorf("User %v doesn't exist", email)
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

	err = s.UserAliasDeleteAll(email)

	return err
}
