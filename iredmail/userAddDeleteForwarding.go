package iredmail

import (
	"fmt"
)

func (s *Server) UserAddForwarding(userEmail, destinationEmail string) (Forwarding, error) {
	f := Forwarding{
		Address:    userEmail,
		Forwarding: destinationEmail,
	}

	userExists, err := s.userExists(userEmail)
	if err != nil {
		return f, err
	}

	if !userExists {
		return f, fmt.Errorf("User %v doesn't exist", userEmail)
	}

	_, userDomain := parseEmail(userEmail)
	_, destDomain := parseEmail(destinationEmail)

	_, err = s.DB.Exec(`
	INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_forwarding)
    VALUES ('` + userEmail + `', '` + destinationEmail + `','` + userDomain + `', '` + destDomain + `', 1);
	`)

	return f, err
}

func (s *Server) UserDeleteForwarding(userAddress, destinationAddress string) error {
	userExists, err := s.userExists(userAddress)
	if err != nil {
		return err
	}
	if !userExists {
		return fmt.Errorf("User %v doesn't exist", userAddress)
	}

	_, err = s.DB.Exec(`DELETE FROM forwardings WHERE address='` + userAddress + `' AND forwarding='` + destinationAddress + `' AND is_forwarding=1;`)

	return err
}
