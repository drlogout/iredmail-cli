package iredmail

import (
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

func (a Forwardings) FilterBy(filter string) Forwardings {
	filteredForwardings := Forwardings{}

	for _, al := range a {
		if strings.Contains(al.Address, filter) {
			filteredForwardings = append(filteredForwardings, al)
		}
	}

	return filteredForwardings
}

func (f Forwardings) GetByAddress(address string) Forwardings {
	filteredForwardings := Forwardings{}

	for _, forwarding := range f {
		if forwarding.Address == address {
			filteredForwardings = append(filteredForwardings, forwarding)
		}
	}

	return filteredForwardings
}

func (s *Server) queryForwardings(query string) (Forwardings, error) {
	Forwardings := Forwardings{}
	rows, err := s.DB.Query(query)
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

func (s *Server) ForwardingList() (Forwardings, error) {
	return s.queryForwardings(`SELECT address, domain, forwarding, dest_domain, active FROM forwardings ORDER BY domain ASC, address ASC;`)
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
