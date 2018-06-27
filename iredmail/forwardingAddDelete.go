package iredmail

import (
	"fmt"
)

func (s *Server) ForwardingAdd(userEmail, destinationEmail string) error {
	userExists, err := s.userExists(userEmail)
	if err != nil {
		return err
	}

	if !userExists {
		return fmt.Errorf("User %v doesn't exist", userEmail)
	}

	_, userDomain := parseEmail(userEmail)
	_, destDomain := parseEmail(destinationEmail)

	_, err = s.DB.Exec(`
	INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_forwarding)
    VALUES ('` + userEmail + `', '` + destinationEmail + `','` + userDomain + `', '` + destDomain + `', 1);
	`)

	return err
}

func (s *Server) ForwardingDelete(userAddress, destinationAddress string) error {
	exists, err := s.forwardingExists(userAddress, destinationAddress)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Forwarding %v -> %v dosn't exist", userAddress, destinationAddress)
	}

	_, err = s.DB.Exec(`DELETE FROM forwardings WHERE address='` + userAddress + `' AND forwarding='` + destinationAddress + `' AND is_forwarding=1;`)

	return err
}

func (s *Server) ForwardingDeleteAll(userAddress string) error {
	_, err := s.DB.Exec(`DELETE FROM forwardings WHERE address='` + userAddress + `' AND is_forwarding=1;`)

	return err
}
