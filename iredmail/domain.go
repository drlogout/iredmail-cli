package iredmail

import (
	"fmt"
	"strconv"
)

func (s *Server) DomainExists(domain string) (bool, error) {
	var exists bool

	query := `select exists
	(select domain from domain
	where domain = '` + domain + `');`

	err := s.DB.QueryRow(query).Scan(&exists)
	if err != nil {
		return exists, err
	}

	if exists {
		return true, nil
	}

	return exists, nil
}

func (s *Server) DomainCreate(domain string, quota int) error {
	settings := fmt.Sprintf("default_user_quota:%v", strconv.Itoa(quota))

	_, err := s.DB.Exec(`
		REPLACE INTO domain (domain, description, settings)
		VALUES ('` + domain + `', '` + domain + `', '` + settings + `')
	`)

	return err
}
