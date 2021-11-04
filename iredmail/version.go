package iredmail

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/blang/semver"
)

const (
	// Version of iredmail-cli
	Version = "0.2.91"

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

	re := regexp.MustCompile(`^\d\.\d\.\d\s*(MYSQL|MARIADB)\s*edition`)
	versionLine := re.FindString(string(file))

	if versionLine == "" {
		return version, fmt.Errorf("No MYSQL nor MariaDB version info found in release file %s", releaseFile)
	}

	splitLine := strings.Split(versionLine, " ")
	version = iredMailVersion(splitLine[0])

	return version, nil
}

// Check checks the iredMail version
func (v *iredMailVersion) Check() error {
	version, err := GetIredMailVersion()
	if err != nil {
		return err
	}

	versionMin, err := semver.Parse(supportedReleaseMin)
	if err != nil {
		return err
	}
	versionMax, err := semver.Parse(supportedReleaseMax)
	if err != nil {
		return err
	}
	versionCur, err := semver.Parse(string(version))
	if err != nil {
		return err
	}

	if versionCur.LT(versionMin) || versionCur.GT(versionMax) {
		return ErrIredMailVersionNotSupported
	}

	return nil
}
