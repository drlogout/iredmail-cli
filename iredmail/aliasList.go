package iredmail

import "github.com/davecgh/go-spew/spew"

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

	spew.Dump(forwardings)

	for i, a := range aliases {
		for _, f := range forwardings {
			if f.Address == a.Address {
				aliases[i].Forwardings = append(a.Forwardings, f)
			}
		}
	}

	return aliases, nil
}
