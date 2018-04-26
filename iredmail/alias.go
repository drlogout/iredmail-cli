package iredmail

import (
	"fmt"
	"sort"
	"strings"
)

type Alias struct {
	Email  string
	Name   string
	Domain string
	Active bool
}

type Aliases []Alias

func (a Aliases) Len() int      { return len(a) }
func (a Aliases) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Aliases) Less(i, j int) bool {
	if a[i].Domain == a[j].Domain {
		usernameSlice := []string{a[i].Name, a[j].Name}
		sort.Strings(usernameSlice)
		if a[i].Name == usernameSlice[0] {
			return true
		}

		return false
	}

	domainSlice := []string{a[i].Domain, a[j].Domain}
	sort.Strings(domainSlice)
	if a[i].Domain == domainSlice[0] {
		return true
	}

	return false
}

func (a Aliases) FilterBy(filter string) Aliases {
	filteredAliases := Aliases{}

	for _, al := range a {
		if strings.Contains(al.Email, filter) {
			filteredAliases = append(filteredAliases, al)
		}
	}

	return filteredAliases
}

func (s *Server) AliasList() (Aliases, error) {
	aliases := Aliases{}
	rows, err := s.DB.Query(`SELECT address, domain, active FROM alias;`)
	if err != nil {
		return aliases, err
	}
	defer rows.Close()

	for rows.Next() {
		var address, domain string
		var active bool

		err := rows.Scan(&address, &domain, &active)
		if err != nil {
			return aliases, err
		}

		name, _ := parseEmail(address)

		aliases = append(aliases, Alias{
			Email:  address,
			Name:   name,
			Domain: domain,
			Active: active,
		})
	}
	err = rows.Err()

	return aliases, err
}

func (s *Server) AliasAdd(email string) error {
	_, domain := parseEmail(email)

	domainExists, err := s.DomainExists(domain)
	if err != nil {
		return err
	}
	if !domainExists {
		return fmt.Errorf("Domain %v does not exist", domain)
	}

	mailboxExists, err := s.MailboxExists(email)
	if err != nil {
		return err
	}
	if !mailboxExists {
		return fmt.Errorf("A mailbox with %v already exists", email)
	}

	_, err = s.DB.Exec(`
		REPLACE INTO alias (address, domain, active)
		VALUES ('` + email + `', '` + domain + `', 1)
	`)

	return err
}
