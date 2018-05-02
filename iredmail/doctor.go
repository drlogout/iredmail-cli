package iredmail

func (s *Server) Doctor() error {

	return s.aliasCheck()
}
