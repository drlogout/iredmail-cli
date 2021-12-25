package integrationTest

import (
	"database/sql"
	"fmt"
	"os/exec"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("mailbox alias add/delete", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can add a mailbox alias", func() {
		if skipMailboxAliasAddDelete && !isCI {
			Skip("can add a mailbox alias")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "mailbox", "add-alias", mailboxAlias1, mailboxName1)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := fmt.Sprintf("Successfully added mailbox alias %s %s %s\n", mailboxAlias1, arrowRight, mailboxName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool
		aliasEmail := mailboxAlias1

		sqlQuery := `SELECT exists
		(SELECT address FROM forwardings
		WHERE address = ? AND forwarding = ?
		AND is_forwarding = 0 AND is_alias = 1 AND is_list = 0 AND active = 1);`

		err = db.QueryRow(sqlQuery, aliasEmail, mailboxName1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))
	})

	It("can't add a mailbox alias if mailbox exists", func() {
		if skipMailboxAliasAddDelete && !isCI {
			Skip("can't add a mailbox alias if mailbox exists")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "mailbox", "add-alias", mailboxName1, mailboxName1)
		output, err = cli.CombinedOutput()
		if err == nil {
			Fail("Expect an error")
		}

		actual := string(output)
		expected := fmt.Sprintf("A mailbox with %s already exists\n", mailboxName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can't add a mailbox alias if mailbox alias already exists", func() {
		if skipMailboxAliasAddDelete && !isCI {
			Skip("can't add a mailbox alias if mailbox alias already exists")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}
		cli = exec.Command(cliPath, "mailbox", "add-alias", mailboxAlias1, mailboxName1)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "mailbox", "add-alias", mailboxAlias1, mailboxName1)
		output, err = cli.CombinedOutput()
		if err == nil {
			Fail("Expect an error")
		}

		actual := string(output)
		expected := fmt.Sprintf("A mailbox alias with %s already exists\n", mailboxAlias1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can't add a mailbox alias if alias already exists", func() {
		if skipMailboxAliasAddDelete && !isCI {
			Skip("can't add a mailbox alias if alias already exists")
		}

		cli := exec.Command(cliPath, "alias", "add", mailboxAlias1)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "mailbox", "add-alias", mailboxAlias1, mailboxName1)
		output, err = cli.CombinedOutput()
		if err == nil {
			Fail("Expect an error")
		}

		actual := string(output)
		expected := fmt.Sprintf("An alias with %s already exists\n", mailboxAlias1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can't add a mailbox alias if mailbox doesn't exist", func() {
		if skipMailboxAliasAddDelete && !isCI {
			Skip("can't add a mailbox alias if mailbox doesn't exist")
		}

		cli := exec.Command(cliPath, "mailbox", "add-alias", mailboxAlias1, mailboxName1)
		output, err := cli.CombinedOutput()
		Expect(err).To(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Mailbox %s doesn't exist\n", mailboxName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can delete a mailbox alias", func() {
		if skipMailboxAliasAddDelete && !isCI {
			Skip("can delete a mailbox alias")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "mailbox", "add-alias", mailboxAlias1, mailboxName1)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "mailbox", "delete-alias", mailboxAlias1)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := fmt.Sprintf("Successfully deleted mailbox alias %s\n", mailboxAlias1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can't delete an alias which doesn't exist", func() {
		if skipMailboxAliasAddDelete && !isCI {
			Skip("can't delete an alias which doesn't exist")
		}

		cli := exec.Command(cliPath, "mailbox", "delete-alias", mailboxName1)
		output, err := cli.CombinedOutput()
		if err == nil {
			Fail("Expect an error")
		}

		actual := string(output)
		expected := fmt.Sprintf("An alias with %s doesn't exists\n", mailboxName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})
})
