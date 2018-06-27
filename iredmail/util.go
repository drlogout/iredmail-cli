package iredmail

import (
	"os/exec"
	"strings"
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
