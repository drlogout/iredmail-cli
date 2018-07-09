package integrationTest

import (
	"database/sql"
	"fmt"
	"os/exec"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("forwarding add/delete", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can add a forwarding", func() {
		if skipForwardingAddDelete && !isCI {
			Skip("can add a forwarding")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "forwarding", "add", mailboxName1, forwardingAddress1)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := fmt.Sprintf("Successfully added forwarding %s %s %s\n", mailboxName1, arrowRight, forwardingAddress1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		sqlQuery := `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = ? AND forwarding = ?
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_list = 0);`

		err = db.QueryRow(sqlQuery, mailboxName1, forwardingAddress1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())

		Expect(exists).To(Equal(true))
	})

	It("can delete a forwarding", func() {
		if skipForwardingAddDelete && !isCI {
			Skip("can delete a forwarding")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "forwarding", "add", mailboxName1, forwardingAddress1)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "forwarding", "delete", mailboxName1, forwardingAddress1)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := fmt.Sprintf("Successfully deleted forwarding %s %s %s\n", mailboxName1, arrowRight, forwardingAddress1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		sqlQuery := `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = ? AND forwarding = ?
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_list = 0);`
		_, err = db.Exec(sqlQuery, mailboxName1, forwardingAddress1)
		Expect(err).NotTo(HaveOccurred())

		Expect(exists).To(Equal(false))
	})

	It("can't add an existing forwarding", func() {
		if skipForwardingAddDelete && !isCI {
			Skip("can't add an existing forwarding")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "forwarding", "add", mailboxName1, forwardingAddress1)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "forwarding", "add", mailboxName1, forwardingAddress1)
		output, err = cli.CombinedOutput()
		if err == nil {
			Fail("Expect an error")
		}

		actual := string(output)
		expected := fmt.Sprintf("Forwarding %s %s %s already exists\n", mailboxName1, arrowRight, forwardingAddress1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can't delete a non-existing forwarding", func() {
		if skipForwardingAddDelete && !isCI {
			Skip("can't delete a non-existing forwarding")
		}

		cli := exec.Command(cliPath, "forwarding", "delete", mailboxName1, forwardingAddress1)
		output, err := cli.CombinedOutput()
		if err == nil {
			Fail("Expect an error")
		}

		actual := string(output)
		expected := fmt.Sprintf("Forwarding %s %s %s doesn't exist\n", mailboxName1, arrowRight, forwardingAddress1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})
})
