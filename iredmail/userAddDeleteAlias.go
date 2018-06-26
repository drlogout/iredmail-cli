package iredmail

import (
	"fmt"
)

func (s *Server) UserAddAlias(alias, email string) error {
	_, domain := parseEmail(email)
	a := fmt.Sprintf("%v@%v", alias, domain)

	userExists, err := s.userExists(a)
	if err != nil {
		return err
	}
	if userExists {
		return fmt.Errorf("An user with %v already exists", a)
	}

	aliasExists, err := s.aliasExists(a)
	if err != nil {
		return err
	}
	if aliasExists {
		return fmt.Errorf("An alias with %v already exists", a)
	}

	_, err = s.DB.Exec(`
		INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_alias, active)
		VALUES ('` + a + `', '` + email + `', '` + domain + `', '` + domain + `', 1, 1)
	`)

	return err
}

func (s *Server) UserDeleteAlias(alias, email string) error {
	_, domain := parseEmail(email)
	a := fmt.Sprintf("%v@%v", alias, domain)

	userExists, err := s.userExists(a)
	if err != nil {
		return err
	}
	if userExists {
		return fmt.Errorf("An user with %v already exists", a)
	}

	aliasExists, err := s.aliasExists(a)
	if err != nil {
		return err
	}
	if aliasExists {
		return fmt.Errorf("An alias with %v already exists", a)
	}

	_, err = s.DB.Exec(`
		INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_alias, active)
		VALUES ('` + a + `', '` + email + `', '` + domain + `', '` + domain + `', 1, 1)
	`)

	return err
}
