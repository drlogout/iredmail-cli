package iredmail

import (
	"sort"
	"strings"
)

type Mailboxes []string

func (m Mailboxes) Len() int      { return len(m) }
func (m Mailboxes) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m Mailboxes) Less(i, j int) bool {
	iSplit := strings.Split(m[i], "@")
	iUsername := iSplit[0]
	iDomain := iSplit[1]

	jSplit := strings.Split(m[j], "@")
	jUsername := jSplit[0]
	jDomain := jSplit[1]

	if iDomain == jDomain {
		usernameSlice := []string{iUsername, jUsername}
		sort.Strings(usernameSlice)
		if iUsername == usernameSlice[0] {
			return true
		}

		return false
	}

	domainSlice := []string{iDomain, jDomain}
	sort.Strings(domainSlice)
	if domainSlice[0] == iDomain {
		return true
	}

	return false
}

func (m Mailboxes) FilterBy(filter string) Mailboxes {
	filteredMailboxes := Mailboxes{}

	for _, mailbox := range m {
		if strings.Contains(mailbox, filter) {
			filteredMailboxes = append(filteredMailboxes, mailbox)
		}
	}

	return filteredMailboxes
}

func (s *Server) MailboxList() (Mailboxes, error) {
	mailboxes := Mailboxes{}
	rows, err := s.DB.Query(`SELECT username FROM mailbox;`)
	if err != nil {
		return mailboxes, err
	}
	defer rows.Close()

	for rows.Next() {
		var username string

		err := rows.Scan(&username)
		if err != nil {
			return mailboxes, err
		}

		mailboxes = append(mailboxes, username)
	}
	err = rows.Err()

	return mailboxes, err
}
