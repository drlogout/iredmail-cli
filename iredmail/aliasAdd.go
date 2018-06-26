package iredmail

import "fmt"

func (s *Server) AliasAdd(email, destEmail string) error {
	_, domain := parseEmail(email)
	_, destDomain := parseEmail(destEmail)

	domainExists, err := s.DomainExists(domain)
	if err != nil {
		return err
	}
	if !domainExists {
		return fmt.Errorf("Domain %v does not exist, please create one first", domain)
	}

	userExists, err := s.userExists(email)
	if err != nil {
		return err
	}
	if userExists {
		return fmt.Errorf("There is already a user %v", email)
	}

	isMailboxAlias, err := s.isUserAlias(email)
	if err != nil {
		return err
	}
	if isMailboxAlias {
		return fmt.Errorf("%v is an alias user", email)
	}

	isAlias, err := s.isAlias(email)
	if err != nil {
		return err
	}
	if !isAlias {
		_, err = s.DB.Exec(`
			INSERT INTO alias (address, domain, active)
			VALUES ('` + email + `', '` + domain + `', 1)
		`)
		if err != nil {
			return err
		}
	}

	_, err = s.DB.Exec(`
		INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_list, active)
		VALUES ('` + email + `', '` + destEmail + `', '` + domain + `', '` + destDomain + `', 1, 1)
	`)

	return err
}
