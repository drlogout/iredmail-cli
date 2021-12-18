package iredmail

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

const (
	// Version of iredmail-cli
	Version = "0.3.0"

	releaseFile         = "/etc/iredmail-release"
	supportedReleaseMin = "0.9.8"
	supportedReleaseMax = "1.3.2"
)

var (
	// ErrIredMailVersionNotSupported ...
	ErrIredMailVersionNotSupported = errors.New("iredMail version is not supported")
)

type iredMailVersion string

// GetIredMailVersion retrievs the iredMail version
func GetIredMailVersion() (iredMailVersion, error) {
	var version iredMailVersion

	if _, err := os.Stat(releaseFile); os.IsNotExist(err) {
		return version, fmt.Errorf("iredMail release file %s does not exist, is iredMail installed?", releaseFile)
	}

	file, err := ioutil.ReadFile(releaseFile)
	if err != nil {
		return version, err
	}

	re := regexp.MustCompile(`(?:^\d\.\d\.\d\s*(MYSQL|MARIADB)\s*edition)|(?:^\d{10} \(Backend: (mariadb|mysql).*)`)
	versionLine := re.FindAllString(string(file), 2)

	if len(versionLine) < 1 {
		return version, fmt.Errorf("No no MYSQL nor MariaDB version info found in release file %s", releaseFile)
	}

	version = iredMailVersion(versionLine[0] + versionLine[1])

	return version, nil
}

// Check checks the iredMail version
func (v *iredMailVersion) Check() error {

	return nil
}
