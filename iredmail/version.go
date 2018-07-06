package iredmail

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/blang/semver"
)

const (
	releaseFile         = "/etc/iredmail-release"
	supportedReleaseMin = "0.9.8"
	supportedReleaseMax = "0.9.8"
)

func Version() (string, error) {
	var version string

	if _, err := os.Stat(releaseFile); os.IsNotExist(err) {
		return version, fmt.Errorf("iredMail release file %s does not exist, is iredMail installed?", releaseFile)
	}

	file, err := ioutil.ReadFile(releaseFile)
	if err != nil {
		return version, err
	}

	re := regexp.MustCompile(`^\d\.\d\.\d\s*MYSQL\s*edition`)
	versionLine := re.FindString(string(file))

	if versionLine == "" {
		return version, fmt.Errorf("No version info found in release file %s", releaseFile)
	}

	splitLine := strings.Split(versionLine, " ")
	version = splitLine[0]

	versionMin, err := semver.Parse(supportedReleaseMin)
	if err != nil {
		return version, err
	}
	versionMax, err := semver.Parse(supportedReleaseMax)
	if err != nil {
		return version, err
	}
	versionCur, err := semver.Parse(version)
	if err != nil {
		return version, err
	}

	if versionCur.LT(versionMin) || versionCur.GT(versionMax) {
		return version, fmt.Errorf("iredMail version %s is not supported", version)
	}

	return version, nil
}
