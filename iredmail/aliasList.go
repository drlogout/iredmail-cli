package iredmail

func (s *Server) AliasList() (Aliases, error) {
	return s.queryAliases(`SELECT address, domain, active FROM alias ORDER BY domain ASC, address ASC;`)
}
