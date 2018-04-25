package iredmail

import (
	"fmt"
	"sort"
	"strings"
)

type Alias struct {
	Address string
	Domain  string
	Active  bool
}

type Aliases []Alias

func (a Aliases) Len() int      { return len(a) }
func (a Aliases) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Aliases) Less(i, j int) bool {
	iName, _ := parseEmail(a[i].Address)
	jName, _ := parseEmail(a[j].Address)

	if a[i].Domain == a[j].Domain {
		usernameSlice := []string{iName, jName}
		sort.Strings(usernameSlice)
		if iName == usernameSlice[0] {
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
		if strings.Contains(al.Address, filter) {
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

		aliases = append(aliases, Alias{
			Address: address,
			Domain:  domain,
			Active:  active,
		})
	}
	err = rows.Err()

	return aliases, err
}

func (s *Server) AliasAdd(email string) error {
	_, domain := parseEmail(email)

	exists, err := s.DomainExists(domain)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("Domain %v does not exist", domain)
	}

	_, err = s.DB.Exec(`
		REPLACE INTO alias (address, domain, active)
		VALUES ('` + email + `', '` + domain + `', 1)
	`)

	return err
}
