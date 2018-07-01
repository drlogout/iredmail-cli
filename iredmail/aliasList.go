package iredmail

import "fmt"

func (s *Server) Aliases() (Aliases, error) {
	aliases, err := s.queryAliases(queryOptions{})
	if err != nil {
		return aliases, err
	}

	forwardings, err := s.queryForwardings(queryOptions{
		where: "is_list = 1",
	})
	if err != nil {
		return aliases, err
	}

	for i, a := range aliases {
		for _, f := range forwardings {
			if f.Address == a.Address {
				aliases[i].Forwardings = append(aliases[i].Forwardings, f)
			}
		}
	}

	return aliases, nil
}

func (s *Server) Alias(aliasEmail string) (Alias, error) {
	alias := Alias{}

	aliasExists, err := s.regularAliasExists(aliasEmail)
	if err != nil {
		return alias, err
	}
	if !aliasExists {
		return alias, fmt.Errorf("Alias %v doesn't exist", aliasEmail)
	}

	aliases, err := s.queryAliases(queryOptions{
		where: "address = '" + aliasEmail + "'",
	})
	if err != nil {
		return alias, err
	}

	if len(aliases) == 0 {
		return alias, fmt.Errorf("Alias not found")
	}

	alias = aliases[0]

	forwardings, err := s.queryForwardings(queryOptions{
		where: "address = '" + aliasEmail + "' AND is_list = 1",
	})
	if err != nil {
		return alias, err
	}

	alias.Forwardings = forwardings

	return alias, nil
}
