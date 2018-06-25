package iredmail

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
)

// types
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

type Users []User

type UserAlias struct {
	Address string
	User    string
}

type UserAliases []UserAlias

func (u User) String() string {
	var buf bytes.Buffer
	w := new(tabwriter.Writer)

	w.Init(&buf, 40, 8, 0, ' ', 0)
	fmt.Fprintf(w, "%v\t%v\n", u.Email, u.Quota)

	w.Flush()

	return buf.String()
}

func (users Users) String() string {
	var buf bytes.Buffer
	w := new(tabwriter.Writer)

	w.Init(&buf, 40, 8, 0, ' ', 0)
	for _, u := range users {
		fmt.Fprintf(w, u.String())
	}

	w.Flush()

	return buf.String()
}

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
	users := Users{}

	whereOption := ""
	if len(options.where) > 1 {
		whereOption = fmt.Sprintf("WHERE %v", options.where)
	}

	rows, err := s.DB.Query(`SELECT username, password, name, domain, quota, maildir FROM mailbox
` + whereOption + `
ORDER BY domain ASC, name ASC;`)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var username, password, name, domain, maildir string
		var quota int

		err := rows.Scan(&username, &password, &name, &domain, &quota, &maildir)
		if err != nil {
			return users, err
		}

		forwardings, err := s.queryForwardings(queryOptions{
			where: "address = '" + username + "'",
		})
		if err != nil {
			return users, err
		}

		users = append(users, User{
			Email:        username,
			Name:         name,
			Domain:       domain,
			PasswordHash: password,
			Quota:        quota,
			MailDir:      maildir,
			Forwardings:  forwardings,
		})
	}
	err = rows.Err()

	return users, err
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

func (s *Server) UserList(args ...string) (Users, error) {
	return s.userQuery(queryOptions{})
}
