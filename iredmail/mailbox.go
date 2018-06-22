package iredmail

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
)

type Mailbox struct {
	Email        string
	Name         string
	Domain       string
	PasswordHash string
	Quota        int
	Type         string
	MailDir      string
	MailboxAliases
	Forwardings
}

func (m Mailbox) Print() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 16, 8, 0, '\t', 0)
	fmt.Fprintf(w, "%v\t%v\n", "Mailbox", "Quota (KB)")
	fmt.Fprintf(w, "%v\t%v\n", "-------", "----------")
	fmt.Fprintf(w, "%v\t%v\n", m.Email, m.Quota)
	w.Flush()
}

type Mailboxes []Mailbox

func (mb Mailboxes) Print() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 16, 8, 0, '\t', 0)
	fmt.Fprintf(w, "%v\t%v\n", "Mailbox", "Quota (KB)")
	fmt.Fprintf(w, "%v\t%v\n", "-------", "----------")
	for _, m := range mb {
		fmt.Fprintf(w, "%v\t%v\n", m.Email, m.Quota)
	}
	w.Flush()
}

type MailboxAlias struct {
	Address string
	Mailbox string
}

type MailboxAliases []MailboxAlias

func (m Mailboxes) FilterBy(filter string) Mailboxes {
	filteredMailboxes := Mailboxes{}

	for _, mailbox := range m {
		if strings.Contains(mailbox.Email, filter) {
			filteredMailboxes = append(filteredMailboxes, mailbox)
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

		mailboxes = append(mailboxes, Mailbox{
			Email:        username,
			Name:         name,
			Domain:       domain,
			PasswordHash: password,
			Quota:        quota,
			MailDir:      maildir,
		})
	}
	err = rows.Err()

	return mailboxes, err
}

func (s *Server) MailboxList() (Mailboxes, error) {
	mailboxes, err := s.mailboxQuery(queryOptions{})

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
		return m, fmt.Errorf("A mailbox %v already exists", email)
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
			storagebasedirectory,storagenode, maildir,
			quota, domain, active, passwordlastchange, created)
		VALUES ('` + email + `', '` + hash + `', '` + name + `',
			'` + storageBase + `','` + storageNode + `', '` + mailDirHash + `',
			'` + strconv.Itoa(quota) + `', '` + domain + `', '1', NOW(), NOW());
		`)
	if err != nil {
		return m, err
	}

	err = s.MailboxAddForwarding(email, email)

	return m, err
}

func (s *Server) mailboxExists(email string) (bool, error) {
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

func (s *Server) MailboxAddAlias(alias, email string) error {
	_, domain := parseEmail(email)
	a := fmt.Sprintf("%v@%v", alias, domain)

	exists, err := s.mailboxExists(a)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("A mailbox with %v already exists", a)
	}

	aliasExists, err := s.aliasExists(a)
	if err != nil {
		return err
	}
	if aliasExists {
		return fmt.Errorf("An alias with %v already exists", a)
	}

	_, err = s.DB.Exec(`
		INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_alias, active)
		VALUES ('` + a + `', '` + email + `', '` + domain + `', '` + domain + `', 1, 1)
	`)

	return err
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

	_, err = s.DB.Exec(`DELETE FROM forwardings WHERE address='` + email + `' AND forwarding='` + email + `' AND is_forwarding=1;`)

	return err
}

func (s *Server) MailboxAddForwarding(mailboxAddress, destinationAddress string) error {
	mailboxExists, err := s.mailboxExists(mailboxAddress)
	if err != nil {
		return err
	}

	if !mailboxExists {
		return fmt.Errorf("Mailbox %v doesn't exist", mailboxAddress)
	}

	_, mailboxDomain := parseEmail(mailboxAddress)
	_, destDomain := parseEmail(destinationAddress)

	_, err = s.DB.Exec(`
	INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_forwarding)
    VALUES ('` + mailboxAddress + `', '` + destinationAddress + `','` + mailboxDomain + `', '` + destDomain + `', 1);
	`)

	return err
}

func (s *Server) MailboxDeleteForwarding(mailbox, destinationAddress string) error {
	mailboxExists, err := s.mailboxExists(mailbox)
	if err != nil {
		return err
	}
	if !mailboxExists {
		return fmt.Errorf("Mailbox %v doesn't exist", mailbox)
	}

	_, err = s.DB.Exec(`DELETE FROM forwardings WHERE address='` + mailbox + `' AND forwarding='` + destinationAddress + `' AND is_forwarding=1;`)

	return err
}
