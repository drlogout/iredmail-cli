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

	hash = strings.TrimSuffix(string(out), "\n")
	return hash, err
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

func PrintAliases(aliases Aliases) {
	var lastAliasDomain string
	var lastAlias string

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 20, 8, 0, '\t', 0)
	fmt.Fprintf(w, "%v\t\t%v\t%v\n", "Alias", "Forwarding", "Type")
	fmt.Fprintf(w, "%v\t\t%v\t%v\n", "-----", "----------", "----")
	for _, a := range aliases {
		if lastAliasDomain != "" && lastAliasDomain != a.Domain {
			fmt.Fprintf(w, "\t\t\t\n")
		}
		lastAliasDomain = a.Domain

		email := a.Address
		arrow := " ->"
		if lastAliasDomain != "" && lastAlias == a.Address {
			email = ""
			arrow = "|->"
		}
		lastAlias = a.Address

		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", email, arrow, "", a.Type)
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
