package integrationTest

import (
	"database/sql"
	"fmt"
	"os/exec"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("domain add/delete", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can add a domain", func() {
		if skipDomainAddDelete && !isCI {
			Skip("can add a domain")
		}

		cli := exec.Command(cliPath, "domain", "add", domain1)
		output, err := cli.CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Successfully added domain %s\n", domain1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT domain FROM domain
		WHERE domain = ?);`

		err = db.QueryRow(query, domain1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))
	})

	It("can't add an existing domain", func() {
		if skipDomainAddDelete && !isCI {
			Skip("can't add an existing domain")
		}

		cli := exec.Command(cliPath, "domain", "add", domain1)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "add", domain1)
		output, err := cli.CombinedOutput()
		Expect(err).To(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Domain %s already exists\n", domain1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can create a domain automatically while adding a mailbox", func() {
		if skipDomainAddDelete && !isCI {
			Skip("can create a domain automatically while adding a mailbox")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName3, mailboxPW)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT domain FROM domain
		WHERE domain = ?);`

		err = db.QueryRow(query, domain1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))
	})

	It("can create a domain automatically while adding an alias", func() {
		if skipDomainAddDelete && !isCI {
			Skip("can create a domain automatically while adding an alias")
		}

		cli := exec.Command(cliPath, "alias", "add", alias1)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT domain FROM domain
		WHERE domain = ?);`

		err = db.QueryRow(query, domain1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))
	})

	It("can delete a domain", func() {
		if skipDomainAddDelete && !isCI {
			Skip("can delete a domain")
		}

		cli := exec.Command(cliPath, "domain", "add", domain1)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "delete", "--force", domain1)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := fmt.Sprintf("Successfully deleted domain %s\n", domain1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT domain FROM domain
		WHERE domain = '` + domain1 + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(false))
	})

	It("can't delete a domain with mailboxes", func() {
		if skipDomainAddDelete && !isCI {
			Skip("can't delete a domain with mailboxes")
		}

		cli := exec.Command(cliPath, "domain", "add", domain1)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "mailbox", "add", mailboxName3, mailboxPW)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "delete", "--force", domain1)
		output, err := cli.CombinedOutput()
		if err == nil {
			Fail("It should exit because mailbox exists")
		}

		actual := string(output)
		expected := fmt.Sprintf("There are still mailboxes with the domain %s, you need to delete them before\n", domain1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT domain FROM domain
		WHERE domain = '` + domain1 + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))
	})

	It("can't delete a domain with aliases", func() {
		if skipDomainAddDelete && !isCI {
			Skip("can't delete a domain with aliases")
		}

		cli := exec.Command(cliPath, "alias", "add", alias1)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "delete", "--force", domain1)
		output, err := cli.CombinedOutput()
		Expect(err).To(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("There are still aliases with the domain %s, you need to delete them before\n", domain1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT domain FROM domain
		WHERE domain = '` + domain1 + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))
	})

	It("can delete a domain with domain aliases", func() {
		if skipDomainAddDelete && !isCI {
			Skip("can delete a domain with domain aliases")
		}

		cli := exec.Command(cliPath, "domain", "add", domain1)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "add-alias", aliasDomain1, domain1)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "add-alias", aliasDomain2, domain1)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "delete", "--force", domain1)
		output, err := cli.CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Successfully deleted domain %s\n", domain1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT domain FROM domain
		WHERE domain = '` + domain1 + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(false))

		query = `SELECT exists
		(SELECT alias_domain FROM alias_domain
		WHERE target_domain = ?);`

		err = db.QueryRow(query, domain1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(false))
	})

	It("can delete a domain with domain catch-all forwardings", func() {
		if skipDomainAddDelete && !isCI {
			Skip("can delete a domain with domain catch-all forwardings")
		}

		cli := exec.Command(cliPath, "domain", "add", domain1)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "add-catchall", domain1, mailboxName1)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "delete", "--force", domain1)
		output, err := cli.CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Successfully deleted domain %s\n", domain1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT domain FROM domain
		WHERE domain = '` + domain1 + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(false))

		query = `SELECT exists
		(SELECT address FROM forwardings
		WHERE address = ? AND forwarding = ? AND is_forwarding = 0 AND is_alias = 0 AND is_list = 0);`

		err = db.QueryRow(query, domain1, mailboxName1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(false))
	})

	It("can't delete an non existing domain", func() {
		if skipDomainAddDelete && !isCI {
			Skip("can't delete an non existing domain")
		}

		cli := exec.Command(cliPath, "domain", "delete", "--force", domain1)
		output, err := cli.CombinedOutput()
		Expect(err).To(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Domain %s doesn't exist\n", domain1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})
})
