package iredmail

import (
	"fmt"
	"sort"
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

func (m Mailboxes) Len() int      { return len(m) }
func (m Mailboxes) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m Mailboxes) Less(i, j int) bool {
	if m[i].Domain == m[j].Domain {
		usernameSlice := []string{m[i].Name, m[j].Name}
		sort.Strings(usernameSlice)
		if m[i].Name == usernameSlice[0] {
			return true
		}

		return false
	}

	domainSlice := []string{m[i].Domain, m[j].Domain}
	sort.Strings(domainSlice)
	if m[i].Domain == domainSlice[0] {
		return true
	}

	return false
}

func (m Mailboxes) FilterBy(filter string) Mailboxes {
	filteredMailboxes := Mailboxes{}

	for _, mailbox := range m {
		if strings.Contains(mailbox.Email, filter) {
			filteredMailboxes = append(filteredMailboxes, mailbox)
		}
	}

	return filteredMailboxes
}

func (s *Server) MailboxList() (Mailboxes, error) {
	mailboxes := Mailboxes{}
	rows, err := s.DB.Query(`SELECT username, password, name, domain, quota FROM mailbox;`)
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

func (s *Server) MailboxCreate(email, password string) error {
	fmt.Println("create email")

	return nil
}
