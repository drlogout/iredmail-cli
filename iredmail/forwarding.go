package iredmail

import (
	"fmt"
	"strings"
)

type Forwarding struct {
	Mailbox             string
	Domain              string
	Forwarding          string
	DestDomain          string
	Active              bool
	IsList              bool
	IsAlias             bool
	IsForwarding        bool
	IsCopyLeftInMailbox bool
}

type Forwardings []Forwarding

type ForwardingsInfo struct {
	IsCopyLeftInMailbox bool
}

type ForwardingsInfoMap map[string]*ForwardingsInfo

func (f *Forwarding) Name() string {
	name, _ := parseEmail(f.Mailbox)

	return name
}

func (forwardings Forwardings) Info() ForwardingsInfoMap {
	infoMap := ForwardingsInfoMap{}

	for _, f := range forwardings {
		_, ok := infoMap[f.Mailbox]
		if !ok {
			infoMap[f.Mailbox] = &ForwardingsInfo{
				IsCopyLeftInMailbox: false,
			}
		}

		if f.Mailbox == f.Forwarding {
			infoMap[f.Mailbox].IsCopyLeftInMailbox = true
		}
	}

	return infoMap
}

func (forwardings Forwardings) External() Forwardings {
	external := Forwardings{}
	for _, f := range forwardings {
		if f.Mailbox != f.Forwarding {
			external = append(external, f)
		}
	}
	return external
}

func (forwardings Forwardings) FilterBy(filter string) Forwardings {
	filteredForwardings := Forwardings{}

	for _, f := range forwardings {
		if strings.Contains(f.Mailbox, filter) {
			filteredForwardings = append(filteredForwardings, f)
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
		var mailboxEmail, domain, forwarding, destDomain string
		var active, isAlias, isForwarding, isList bool

		err := rows.Scan(&mailboxEmail, &domain, &forwarding, &destDomain, &active, &isAlias, &isForwarding, &isList)
		if err != nil {
			return Forwardings, err
		}

		Forwardings = append(Forwardings, Forwarding{
			Mailbox:      mailboxEmail,
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
