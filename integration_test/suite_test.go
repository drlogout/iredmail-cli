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
	dbConnectionStringCI    = "vmail:itslocalletmein@tcp(127.0.0.1:3306)/vmail"
)

var (
	cliPath            string
	dbConnectionString = dbConnectionStringLocal
	projectDir         string
	dbTables           = []string{
		"alias",
		"domain",
		"alias_domain",
		"forwardings",
		"mailbox",
	}
	isCI = false

	mailboxName1 = "mail@example.net"
	mailboxName2 = "info@domain.com"
	mailboxName3 = "webmaster@example.com"
	mailboxName4 = "abuse@domain.com"
	mailboxName5 = "support@example.org"

	mailboxPW  = "alskdlqkdjalskd"
	mailboxPW2 = "qweqoiwueoqwiueq"

	forwardingAddress1 = "info@otherdomain.com"
	forwardingAddress2 = "webmaster@otherexample.net"

	mailboxAlias1 = "postmaster"
	mailboxAlias2 = "abuse"
	mailboxAlias3 = "webmaster"

	alias1 = "developer@example.com"
	alias2 = "support@domain.com"
	alias3 = "help@example.net"

	aliasForwarding1 = "mail@example.com"
	aliasForwarding2 = "info@domain.com"
	aliasForwarding3 = "whatever@otherexample.com"
	aliasForwarding4 = "whatever@example.net"

	domain1 = "example.com"
	domain2 = "example.net"
	domain3 = "domain.com"

	aliasDomain1 = "alias.com"
	aliasDomain2 = "alias.de"

	domainSettings    = "default_user_quota:4096"
	domainDescription = "A domain description"

	customQuota       = 4096
	customStoragePath = "/var/mail/custom"

	arrowRight = "➞"

	skipAliasAddDelete           = false
	skipAliasForwardingAddDelete = false
	skipAliasInfo                = false
	skipAliasList                = false
	skipDomainAddDelete          = false
	skipDomainAliasAddDelete     = false
	skipDomainCatchallAddDelete  = false
	skipDomainList               = false
	skipForwardingAddDelete      = false
	skipMailboxAddDelete         = false
	skipMailboxAliasAddDelete    = false
	skipMailboxInfo              = false
	skipMailboxList              = false
	skipMailboxUpdate            = false
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
