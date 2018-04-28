package iredmail

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

func parseEmail(email string) (string, string) {
	split := strings.Split(email, "@")

	if len(split) < 2 {
		return email, ""
	}

	return split[0], split[1]
}

func generatePassword(password string) (string, error) {
	var hash string

	out, err := exec.Command("doveadm", "pw", "-p", password).Output()
	if err != nil {
		return hash, err
	}

	return string(out), err
}

func generateMaildirHash(email string) string {
	var str1, str2, str3 string

	name, domain := parseEmail(email)
	date := time.Now().UTC().Format("2006.01.02.15.04.05")

	switch len(name) {
	case 1:
		str1 = string(name[0])
		str2 = str1
		str3 = str2
	case 2:
		str1 = string(name[0])
		str2 = string(name[1])
		str3 = str2
	default:
		str1 = string(name[0])
		str2 = string(name[1])
		str3 = string(name[2])
	}

	return domain + `/` + str1 + `/` + str2 + `/` + str3 + `/` + name + `-` + date + `/`
}

func PrintMailboxes(mailboxes Mailboxes) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 16, 8, 0, '\t', 0)
	fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", "Email (user name)", "Quota", "Name", "Domain")
	fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", "-----------------", "-----", "----", "------")
	for _, m := range mailboxes {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", m.Email, m.Quota, m.Name, m.Domain)
	}
	w.Flush()
}

func PrintAliases(aliases Aliases) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 16, 8, 0, '\t', 0)
	fmt.Fprintf(w, "%v\t%v\t%v\n", "Alias", "Domain", "Active")
	fmt.Fprintf(w, "%v\t%v\t%v\n", "-----", "------", "------")
	for _, a := range aliases {
		fmt.Fprintf(w, "%v\t%v\t%v\n", a.Name, a.Domain, a.Active)
	}
	w.Flush()
}

func PrintForwardings(forwardings Forwardings) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 16, 8, 0, '\t', 0)
	fmt.Fprintf(w, "%v\t%v\n", "Address", "Forwarding")
	fmt.Fprintf(w, "%v\t%v\n", "-------", "----------")
	for _, f := range forwardings {
		fmt.Fprintf(w, "%v\t%v\n", f.Address, f.Forwarding)
	}
	w.Flush()
}

func PrintDomains(domains Domains, quiet bool) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 16, 8, 0, '\t', 0)
	if quiet {
		for _, d := range domains {
			fmt.Fprintf(w, "%v\n", d.Domain)
		}
	} else {
		w.Init(os.Stdout, 16, 8, 0, '\t', 0)
		fmt.Fprintf(w, "%v\t%v\t%v\n", "Domain", "Description", "Settings")
		fmt.Fprintf(w, "%v\t%v\t%v\n", "------", "-----------", "--------")
		for _, d := range domains {
			fmt.Fprintf(w, "%v\t%v\t%v\n", d.Domain, d.Description, d.Settings)
		}
	}
	w.Flush()
}

func PrintDomainInfo(domain Domain, mailboxes Mailboxes, aliases Aliases, aliasForwardings Forwardings) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 16, 8, 0, '\t', 0)
	fmt.Fprintf(w, "%v\t%v\n", "Domain", "Description")
	fmt.Fprintf(w, "%v\t%v\n", "------", "-----------")
	fmt.Fprintf(w, "%v\t%v\n", domain.Domain, domain.Description)

	fmt.Fprintln(w)
	fmt.Fprintf(w, "%v\t%v\n", "Mailboxes ("+strconv.Itoa(len(mailboxes))+")", "Quota")
	fmt.Fprintf(w, "%v\t%v\n", "---------", "-----")
	for _, m := range mailboxes {
		fmt.Fprintf(w, "%v\t%v\n", m.Email, m.Quota)
	}

	fmt.Fprintln(w)
	fmt.Fprintf(w, "%v\n", "Aliases ("+strconv.Itoa(len(aliases))+")")
	fmt.Fprintf(w, "%v\n", "-------")
	for _, a := range aliases {
		aliasText := a.Email
		if !a.Active {
			aliasText = aliasText + " (inactive)"
		}
		fmt.Fprintf(w, "%v\n", aliasText)
		forwardings := aliasForwardings.GetByAddress(a.Email)
		for _, f := range forwardings {
			forwardingText := f.Forwarding
			if !f.Active {
				forwardingText = forwardingText + " (inactive)"
			}
			fmt.Fprintf(w, "%v\n", "  ->  "+forwardingText)
		}
	}

	w.Flush()
}
