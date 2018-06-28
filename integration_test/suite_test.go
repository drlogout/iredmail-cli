package integrationTest

import (
	"database/sql"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbConnectionStringLocal = "vmail:sx4fDttWdWNbiBPsGxhbbxic2MmmGsmJ@tcp(127.0.0.1:8806)/vmail"
	dbConnectionStringCI    = "vmail:MDmPEwViyNNrMVpxrRGQivvFdtxZAp98@tcp(iredmail-test-db.noltech.net:3357)/vmail"
)

var (
	cliPath            string
	dbConnectionString = dbConnectionStringLocal
	projectDir         string
	dbTables           = []string{
		"alias",
		"domain",
		"forwardings",
		"mailbox",
	}
	isCI = false

	mailboxName1 = "post@web.de"
	mailboxName2 = "info@domain.com"
	mailboxName3 = "webmaster@example.com"
	mailboxName4 = "abuse@domain.com"
	mailboxName5 = "support@wurst.de"

	mailboxPW = "alskdlqkdjalskd"

	forwardingAddress1 = "info@example.com"
	forwardingAddress2 = "webmaster@example.net"

	alias1 = "mail"
	alias2 = "abuse"
	alias3 = "webmaster"

	customQuota       = 4096
	customStoragePath = "/var/mail/custom"

	skipMailbox      = false
	skipMailboxInfo  = false
	skipMailboxList  = false
	skipMailboxAlias = false
	skipForwarding   = false
)

func TestCLI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}

var _ = BeforeSuite(func() {
	isCI = os.Getenv("CI") == "true"
	if isCI {
		dbConnectionString = dbConnectionStringCI
	}

	err := resetDB()
	Expect(err).NotTo(HaveOccurred())

	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	projectDir = filepath.Join(cwd, "../")
	cliPath = filepath.Join(projectDir, "iredmail-cli")

	cmd := exec.Command("go", "build", "-o", cliPath)
	cmd.Dir = projectDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		Fail(string(output))
	}
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

func loadGolden(filename string) string {
	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	content, err := ioutil.ReadFile(filepath.Join(cwd, "golden", filename))
	Expect(err).NotTo(HaveOccurred())

	return string(content)
}
