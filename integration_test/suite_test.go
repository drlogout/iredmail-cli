package integrationTest

import (
	"database/sql"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbConnectionString = "vmail:sx4fDttWdWNbiBPsGxhbbxic2MmmGsmJ@tcp(127.0.0.1:8806)/vmail"
)

var (
	cliPath    string
	projectDir string
	dbTables   = []string{
		"alias",
		"domain",
		"forwardings",
		"mailbox",
	}
)

func TestCLI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}

var _ = BeforeSuite(func() {
	err := resetDB()
	Expect(err).NotTo(HaveOccurred())

	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	projectDir = filepath.Join(cwd, "../")
	cliPath = filepath.Join(projectDir, "iredmail-cli")

	cmd := exec.Command("go", "build", "-o", cliPath)
	cmd.Dir = projectDir

	err = cmd.Run()
	Expect(err).NotTo(HaveOccurred())
})

func resetDB() error {
	db, err := sql.Open("mysql", dbConnectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	for _, table := range dbTables {
		_, err := db.Exec("DELETE FROM " + table)
		if err != nil {
			return err
		}
	}

	return nil
}
