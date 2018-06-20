package iredmail

func (s *Server) AliasList() (Aliases, error) {
	aliases := Aliases{}
	forwardings, err := s.queryForwardings(queryOptions{})
	if err != nil {
		return aliases, err
	}

	for _, f := range forwardings {
		isRegularAlias, err := s.isAlias(f.Address)
		if err != nil {
			return aliases, err
		}
		if isRegularAlias {
			aliases = append(aliases, Alias{
				Address: f.Address,
				Domain:  f.Domain,
			})
		}
	}

	return aliases, nil
}
