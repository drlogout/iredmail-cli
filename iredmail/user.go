package iredmail

import (
	"fmt"
	"strings"
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
	UserAliases  Forwardings
	Forwardings
}

type Users []User

func (users Users) FilterBy(filter string) Users {
	filteredUsers := Users{}

	for _, user := range users {
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
			where: "address = '" + username + "' AND is_forwarding = 1",
		})
		if err != nil {
			return users, err
		}

		userAliases, err := s.queryForwardings(queryOptions{
			where: "forwarding = '" + username + "' AND is_alias = 1",
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
			UserAliases:  userAliases,
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

func (s *Server) Users() (Users, error) {
	return s.userQuery(queryOptions{})
}

func (s *Server) User(email string) (User, error) {
	exists, err := s.userExists(email)
	if err != nil {
		return User{}, err
	}

	if !exists {
		return User{}, fmt.Errorf("User does not exist")
	}

	users, err := s.userQuery(queryOptions{
		where: `username = '` + email + `'`,
	})
	if err != nil {
		return User{}, err
	}
	if len(users) == 0 {
		return User{}, fmt.Errorf("User not found")
	}

	return users[0], nil
}
