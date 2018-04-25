package iredmail

func (s *Server) ForwardingAdd(address, forwarding string) error {
	_, domain := parseEmail(address)
	_, destDomain := parseEmail(forwarding)

	_, err := s.DB.Exec(`
		INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_forwarding)
    VALUES ('` + address + `', '` + forwarding + `','` + domain + `', '` + destDomain + `', 1);
		`)

	return err
}
