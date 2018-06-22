package iredmail

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
)

type User struct {
	Email        string
	Name         string
	Domain       string
	PasswordHash string
	Quota        int
	Type         string
	MailDir      string
	UserAliases
	Forwardings
}

func (m User) Print() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 16, 8, 0, '\t', 0)
	fmt.Fprintf(w, "%v\t%v\n", "User", "Quota (KB)")
	fmt.Fprintf(w, "%v\t%v\n", "----", "----------")
	fmt.Fprintf(w, "%v\t%v\n", m.Email, m.Quota)
	w.Flush()
}

type Users []User

func (mb Users) Print() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 16, 8, 0, '\t', 0)
	fmt.Fprintf(w, "%v\t%v\n", "User", "Quota (KB)")
	fmt.Fprintf(w, "%v\t%v\n", "----", "----------")
	for _, m := range mb {
		fmt.Fprintf(w, "%v\t%v\n", m.Email, m.Quota)
	}
	w.Flush()
}

type UserAlias struct {
	Address string
	User    string
}

type UserAliases []UserAlias

func (m Users) FilterBy(filter string) Users {
	filteredUsers := Users{}

	for _, user := range m {
		if strings.Contains(user.Email, filter) {
			filteredUsers = append(filteredUsers, user)
		}
	}

	return filteredUsers
}

func (s *Server) userQuery(options queryOptions) (Users, error) {
	useres := Users{}

	whereOption := ""
	if len(options.where) > 1 {
		whereOption = fmt.Sprintf("WHERE %v", options.where)
	}

	rows, err := s.DB.Query(`SELECT username, password, name, domain, quota, maildir FROM mailbox
` + whereOption + `
ORDER BY domain ASC, name ASC;`)
	if err != nil {
		return useres, err
	}
	defer rows.Close()

	for rows.Next() {
		var username, password, name, domain, maildir string
		var quota int

		err := rows.Scan(&username, &password, &name, &domain, &quota, &maildir)
		if err != nil {
			return useres, err
		}

		useres = append(useres, User{
			Email:        username,
			Name:         name,
			Domain:       domain,
			PasswordHash: password,
			Quota:        quota,
			MailDir:      maildir,
		})
	}
	err = rows.Err()

	return useres, err
}

func (s *Server) UserList() (Users, error) {
	useres, err := s.userQuery(queryOptions{})

	return useres, err
}

func (s *Server) UserAdd(email, password string, quota int, storageBasePath string) (User, error) {
	name, domain := parseEmail(email)
	m := User{
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

	userExists, err := s.userExists(email)
	if err != nil {
		return m, err
	}
	if userExists {
		return m, fmt.Errorf("User %v already exists", email)
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

	err = s.UserAddForwarding(email, email)

	return m, err
}

func (s *Server) userExists(email string) (bool, error) {
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

func (s *Server) UserAddAlias(alias, email string) error {
	_, domain := parseEmail(email)
	a := fmt.Sprintf("%v@%v", alias, domain)

	exists, err := s.userExists(a)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("A user with %v already exists", a)
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

	_, err = s.DB.Exec(`DELETE FROM forwardings WHERE address='` + email + `' AND forwarding='` + email + `' AND is_forwarding=1;`)

	return err
}

func (s *Server) UserAddForwarding(userAddress, destinationAddress string) error {
	userExists, err := s.userExists(userAddress)
	if err != nil {
		return err
	}

	if !userExists {
		return fmt.Errorf("User %v doesn't exist", userAddress)
	}

	_, userDomain := parseEmail(userAddress)
	_, destDomain := parseEmail(destinationAddress)

	_, err = s.DB.Exec(`
	INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_forwarding)
    VALUES ('` + userAddress + `', '` + destinationAddress + `','` + userDomain + `', '` + destDomain + `', 1);
	`)

	return err
}

func (s *Server) UserDeleteForwarding(userAddress, destinationAddress string) error {
	userExists, err := s.userExists(userAddress)
	if err != nil {
		return err
	}
	if !userExists {
		return fmt.Errorf("User %v doesn't exist", userAddress)
	}

	_, err = s.DB.Exec(`DELETE FROM forwardings WHERE address='` + userAddress + `' AND forwarding='` + destinationAddress + `' AND is_forwarding=1;`)

	return err
}
