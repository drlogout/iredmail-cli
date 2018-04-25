package iredmail

import (
	"sort"
	"strings"
)

type Forwarding struct {
	Address    string
	Domain     string
	Forwarding string
	DestDomain string
	Active     bool
}

type Forwardings []Forwarding

func (a Forwardings) Len() int      { return len(a) }
func (a Forwardings) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Forwardings) Less(i, j int) bool {
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

func (a Forwardings) FilterBy(filter string) Forwardings {
	filteredForwardings := Forwardings{}

	for _, al := range a {
		if strings.Contains(al.Address, filter) {
			filteredForwardings = append(filteredForwardings, al)
		}
	}

	return filteredForwardings
}

func (s *Server) ForwardingList() (Forwardings, error) {
	Forwardings := Forwardings{}
	rows, err := s.DB.Query(`SELECT address, domain, forwarding, dest_domain, active FROM forwardings;`)
	if err != nil {
		return Forwardings, err
	}
	defer rows.Close()

	for rows.Next() {
		var address, domain, forwarding, destDomain string
		var active bool

		err := rows.Scan(&address, &domain, &forwarding, &destDomain, &active)
		if err != nil {
			return Forwardings, err
		}

		Forwardings = append(Forwardings, Forwarding{
			Address:    address,
			Domain:     domain,
			Forwarding: forwarding,
			DestDomain: destDomain,
			Active:     active,
		})
	}
	err = rows.Err()

	return Forwardings, err
}

func (s *Server) ForwardingAdd(address, forwarding string) error {
	_, domain := parseEmail(address)
	_, destDomain := parseEmail(forwarding)

	_, err := s.DB.Exec(`
		INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_forwarding)
    VALUES ('` + address + `', '` + forwarding + `','` + domain + `', '` + destDomain + `', 1);
		`)

	return err
}
