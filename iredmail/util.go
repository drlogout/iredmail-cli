package iredmail

import (
	"fmt"
	"os"
	"os/exec"
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

func PrintDomains(domains Domains) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 16, 8, 0, '\t', 0)
	fmt.Fprintf(w, "%v\t%v\t%v\n", "Domain", "Description", "Settings")
	fmt.Fprintf(w, "%v\t%v\t%v\n", "------", "-----------", "--------")
	for _, d := range domains {
		fmt.Fprintf(w, "%v\t%v\t%v\n", d.Domain, d.Description, d.Settings)
	}
	w.Flush()
}

func PrintDomainInfo(domain Domain, mailboxes Mailboxes, aliases Aliases) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 16, 8, 0, '\t', 0)
	fmt.Fprintf(w, "%v\t%v\n", "Domain", "Description")
	fmt.Fprintf(w, "%v\t%v\n", "------", "-----------")
	fmt.Fprintf(w, "%v\t%v\n", domain.Domain, domain.Description)

	fmt.Fprintln(w)
	fmt.Fprintf(w, "%v\t%v\n", "Mailboxes", "Quota")
	fmt.Fprintf(w, "%v\t%v\n", "---------", "-----")
	for _, m := range mailboxes {
		fmt.Fprintf(w, "%v\t%v\n", m.Email, m.Quota)
	}

	fmt.Fprintln(w)
	fmt.Fprintf(w, "%v\t%v\n", "Aliases", "Active")
	fmt.Fprintf(w, "%v\t%v\n", "-------", "------")
	for _, a := range aliases {
		fmt.Fprintf(w, "%v\t%v\n", a.Email, a.Active)
	}

	w.Flush()
}
