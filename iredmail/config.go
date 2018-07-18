package iredmail

import (
	"bufio"
	"os"
	"strings"
)

var (
	configPath = ""
	config     = map[string]string{}
)

func ReadInConfig() error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") {
			split := strings.Split(line, "=")
			config[split[0]] = strings.Trim(split[1], " ")
			config[split[0]] = strings.Trim(split[1], "\"")
		}
	}

	return nil
}

func SetConfigFile(path string) {
	configPath = path
}
