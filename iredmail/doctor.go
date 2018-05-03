package iredmail

import "fmt"

func (s *Server) Doctor(repair bool) error {
	if repair {
		fmt.Printf("Running doctor (repair mode)\n")
	}
	return s.checkAliases(repair)
}
