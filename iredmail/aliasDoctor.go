package iredmail

import (
	"fmt"
)

func (s *Server) checkAliases(repair bool) error {
	forwardings, err := s.queryForwardings(queryOptions{})
	if err != nil {
		return err
	}

	for _, f := range forwardings {
		isAlias, err := s.isAlias(f.Address)
		if err != nil {
			return err
		}

		if isAlias {
			err := s.repairAliasforwardings(f.Address, repair)
			if err != nil {
				return err
			}
		}

		forwardingIsMailbox, err := s.MailboxExists(f.Forwarding)
		if err != nil {
			return err
		}

		// if f.Domain and f.DestDomain are equal and f.Forwarding is a local mailbox
		// it should a mailbox alias
		if f.Domain == f.DestDomain && forwardingIsMailbox {
			isMailboxAlias, err := s.isMailboxAlias(f.Address)
			if err != nil {
				return err
			}

			if !isMailboxAlias {
				// fmt.Printf("%v should be a mailbox alias\n", f.Forwarding)
			}
		}
	}

	return nil
}

func (s *Server) repairAliasforwardings(aliasAddress string, repair bool) error {
	falseAliasForwardings := Forwardings{}

	aliasForwardings, err := s.queryForwardings(queryOptions{
		where: fmt.Sprintf("address='%v'", aliasAddress),
	})
	if err != nil {
		return err
	}

	for _, af := range aliasForwardings {
		if !af.IsList || af.IsAlias || af.IsForwarding {
			falseAliasForwardings = append(falseAliasForwardings, af)
			fmt.Printf("%v -> %v, IS: is_list: %v | is_forwarding: %v | is_alias: %v SHOULD: is_list: true | is_forwarding: false | is_alias: false\n", af.Address, af.Forwarding, af.IsList, af.IsForwarding, af.IsAlias)
		}
	}

	if repair {
		for _, af := range falseAliasForwardings {
			_, err := s.DB.Exec(`
				REPLACE INTO forwardings (address, forwarding, domain, dest_domain, is_list, is_forwarding, is_alias, is_maillist, active)
				VALUES ('` + af.Address + `', '` + af.Forwarding + `', '` + af.Domain + `', '` + af.DestDomain + `', 1, 0, 0, 0, 1)
			`)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
