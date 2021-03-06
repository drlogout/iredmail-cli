package integrationTest

import (
	"database/sql"
	"fmt"
	"os/exec"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("domain add/delete-catch-all forwarding", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can add an catch-all forwarding", func() {
		if skipDomainCatchallAddDelete && !isCI {
			Skip("can add an catch-all forwarding")
		}

		cli := exec.Command(cliPath, "domain", "add", domain1)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "domain", "add-catchall", domain1, mailboxName1)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := fmt.Sprintf("Successfully added catch-all forwarding %s %s %s\n", domain1, arrowRight, mailboxName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT address FROM forwardings
		WHERE address = ? AND forwarding = ? AND is_forwarding = 0 AND is_alias = 0 AND is_list = 0);`

		err = db.QueryRow(query, domain1, mailboxName1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))
	})

	It("can't add an existing catch-all forwarding", func() {
		if skipDomainCatchallAddDelete && !isCI {
			Skip("can't add an existing catch-all forwarding")
		}

		cli := exec.Command(cliPath, "domain", "add", domain1)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "domain", "add-catchall", domain1, mailboxName1)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "domain", "add-catchall", domain1, mailboxName1)
		output, err = cli.CombinedOutput()
		if err == nil {
			Fail("Expect an error")
		}

		actual := string(output)
		expected := fmt.Sprintf("Catch-all forwarding %s %s %s already exists\n", domain1, arrowRight, mailboxName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can't add an catch-all forwarding if domain doesn't exist", func() {
		if skipDomainCatchallAddDelete && !isCI {
			Skip("can't add an catch-all forwarding if domain doesn't existt")
		}

		cli := exec.Command(cliPath, "domain", "add-catchall", domain1, mailboxName1)
		output, err := cli.CombinedOutput()
		if err == nil {
			Fail("Expect an error")
		}

		actual := string(output)
		expected := fmt.Sprintf("Domain %s doesn't exists\n", domain1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can delete a catch-all forwarding", func() {
		if skipDomainCatchallAddDelete && !isCI {
			Skip("can delete a catch-all forwarding")
		}

		cli := exec.Command(cliPath, "domain", "add", domain1)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "domain", "add-catchall", domain1, mailboxName1)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "domain", "delete-catchall", domain1, mailboxName1)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := fmt.Sprintf("Successfully deleted catch-all forwarding %s %s %s\n", domain1, arrowRight, mailboxName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT address FROM forwardings
		WHERE address = ? AND forwarding = ? AND is_forwarding = 0 AND is_alias = 0 AND is_list = 0);`

		err = db.QueryRow(query, domain1, mailboxName1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(false))
	})

	It("can't delete a non existing catch-all forwarding", func() {
		if skipDomainCatchallAddDelete && !isCI {
			Skip("can't delete a non existing catch-all forwarding")
		}

		cli := exec.Command(cliPath, "domain", "add", domain1)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "domain", "delete-catchall", domain1, mailboxName1)
		output, err = cli.CombinedOutput()
		if err == nil {
			Fail("Expect an error")
		}

		actual := string(output)
		expected := fmt.Sprintf("Catch-all forwarding %s %s %s doesn't exist\n", domain1, arrowRight, mailboxName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})
})
