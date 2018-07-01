package iredmail

import (
	"strings"
)

type Forwarding struct {
	Address             string
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
	name, _ := parseEmail(f.Address)

	return name
}

func (forwardings Forwardings) Info() ForwardingsInfoMap {
	infoMap := ForwardingsInfoMap{}

	for _, f := range forwardings {
		_, ok := infoMap[f.Address]
		if !ok {
			infoMap[f.Address] = &ForwardingsInfo{
				IsCopyLeftInMailbox: false,
			}
		}

		if f.Address == f.Forwarding {
			infoMap[f.Address].IsCopyLeftInMailbox = true
		}
	}

	return infoMap
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

func (forwardings Forwardings) FilterBy(filter string) Forwardings {
	filteredForwardings := Forwardings{}

	for _, f := range forwardings {
		if strings.Contains(f.Address, filter) {
			filteredForwardings = append(filteredForwardings, f)
		}
	}

	return filteredForwardings
}

func (s *Server) queryForwardings(sqlQuery string, args ...interface{}) (Forwardings, error) {
	Forwardings := Forwardings{}

	rows, err := s.DB.Query(sqlQuery, args)
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
			Address:      mailboxEmail,
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

func (s *Server) forwardingExists(mailboxEmail, destinationEmail string) (bool, error) {
	var exists bool

	sqlQuery := `
	SELECT exists
	(SELECT address FROM forwardings
	WHERE address = ? AND forwarding = ? AND is_forwarding = 1
	);`
	err := s.DB.QueryRow(sqlQuery, mailboxEmail, destinationEmail).Scan(&exists)

	return exists, err
}

func (s *Server) ForwardingList() (Forwardings, error) {
	sqlQuery := `
	SELECT address, domain, forwarding, dest_domain, active, is_alias, is_forwarding, is_list 
	FROM forwardings
	WHERE is_forwarding = 1
	ORDER BY domain ASC, address ASC;`

	return s.queryForwardings(sqlQuery)
}
