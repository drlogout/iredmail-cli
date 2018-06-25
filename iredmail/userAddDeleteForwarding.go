package iredmail

import (
	"fmt"
)

func (s *Server) UserAddForwarding(userAddress, destinationAddress string) error {
	userExists, err := s.userExists(userAddress)
	if err != nil {
		return err
	}

	if !userExists {
		return fmt.Errorf("User %v doesn't exist", userAddress)
	}

	_, userDomain := parseEmail(userAddress)
	_, destDomain := parseEmail(destinationAddress)

	_, err = s.DB.Exec(`
	INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_forwarding)
    VALUES ('` + userAddress + `', '` + destinationAddress + `','` + userDomain + `', '` + destDomain + `', 1);
	`)

	return err
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
