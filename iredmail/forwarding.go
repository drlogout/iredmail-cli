package iredmail

import (
	"fmt"
	"strings"
)

type Forwarding struct {
	Address      string
	Domain       string
	Forwarding   string
	DestDomain   string
	Active       bool
	IsList       bool
	IsAlias      bool
	IsForwarding bool
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

func (s *Server) queryForwardings(options queryOptions) (Forwardings, error) {
	Forwardings := Forwardings{}

	whereOption := ""
	if len(options.where) > 1 {
		whereOption = fmt.Sprintf("WHERE %v", options.where)
	}

	rows, err := s.DB.Query(`SELECT address, domain, forwarding, dest_domain, active, is_alias, is_forwarding, is_list FROM forwardings
` + whereOption + `
ORDER BY domain ASC, address ASC;`)
	if err != nil {
		return Forwardings, err
	}
	defer rows.Close()

	for rows.Next() {
		var address, domain, forwarding, destDomain string
		var active, isAlias, isForwarding, isList bool

		err := rows.Scan(&address, &domain, &forwarding, &destDomain, &active, &isAlias, &isForwarding, &isList)
		if err != nil {
			return Forwardings, err
		}

		Forwardings = append(Forwardings, Forwarding{
			Address:      address,
			Domain:       domain,
			Forwarding:   forwarding,
			DestDomain:   destDomain,
			Active:       active,
			IsAlias:      isAlias,
			IsForwarding: isForwarding,
			IsList:       isList,
		})
	}
	err = rows.Err()

	return Forwardings, err
}

func (s *Server) ForwardingList() (Forwardings, error) {
	return s.queryForwardings(queryOptions{
		where: "domain='wirtschaft-symposium.de'",
	})
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
