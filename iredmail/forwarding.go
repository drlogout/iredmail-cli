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

func (f *Forwarding) Name() string {
	name, _ := parseEmail(f.Address)

	return name
}

func (forwardings Forwardings) IsCopyLeftInMailbox() bool {
	for _, f := range forwardings {
		if f.Address == f.Forwarding {
			return true
		}
	}
	return false
}

func (forwardings Forwardings) External() Forwardings {
	external := Forwardings{}
	for _, f := range forwardings {
		if f.Address != f.Forwarding {
			external = append(external, f)
		}
	}
	return external
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

func (s *Server) queryForwardings(options queryOptions) (Forwardings, error) {
	Forwardings := Forwardings{}

	whereOption := ""
	if len(options.where) > 1 {
		whereOption = fmt.Sprintf("WHERE %v", options.where)
	}

	rows, err := s.DB.Query(`
		SELECT address, domain, forwarding, dest_domain, active, is_alias, is_forwarding, is_list 
		FROM forwardings
		` + whereOption + `
		ORDER BY domain ASC, address ASC;
	`)
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

func (s *Server) forwardingExists(userAddress, destinationAddress string) (bool, error) {
	var exists bool

	query := `
		SELECT exists
		(SELECT * FROM forwardings
			WHERE address = '` + userAddress + `' AND forwarding = '` + destinationAddress + `' AND is_forwarding = 1
		);`
	err := s.DB.QueryRow(query).Scan(&exists)
	if err != nil {
		return exists, err
	}

	if exists {
		return true, nil
	}

	return exists, nil
}

func (s *Server) ForwardingList() (Forwardings, error) {
	forwardings, err := s.queryForwardings(queryOptions{
		where: "is_forwarding = 1",
	})
	if err != nil {
		return forwardings, err
	}

	return forwardings, nil
}
