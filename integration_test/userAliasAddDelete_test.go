package integrationTest

import (
	"database/sql"
	"fmt"
	"os/exec"
	"reflect"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("user alias", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can add an user alias", func() {
		if skipUserAlias && !isCI {
			Skip("can add an user alias")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "mailbox", "add-alias", alias1, mailboxName1)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := fmt.Sprintf("Successfully added user alias %v -> %v\n", alias1, mailboxName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool
		domain := strings.Split(mailboxName1, "@")[1]

		query := `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = '` + fmt.Sprintf("%v@%v", alias1, domain) + `' AND forwarding = '` + mailboxName1 + `'
		AND is_alias = 1 AND active = 1 AND is_forwarding = 0 AND is_maillist = 0);`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))
	})

	It("can't add an user alias if email exists", func() {
		if skipUserAlias && !isCI {
			Skip("can't add an user alias if email exists")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		name := strings.Split(mailboxName1, "@")[0]

		cli = exec.Command(cliPath, "mailbox", "add-alias", name, mailboxName1)
		output, err := cli.CombinedOutput()
		Expect(err).To(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("An user with %v already exists\n", mailboxName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can't add an user alias if user alias already exists", func() {
		if skipUserAlias && !isCI {
			Skip("can't add an user alias if user alias already exists")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "mailbox", "add-alias", alias1, mailboxName1)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "mailbox", "add-alias", alias1, mailboxName1)
		output, err := cli.CombinedOutput()
		Expect(err).To(HaveOccurred())

		domain := strings.Split(mailboxName1, "@")[1]

		actual := string(output)
		expected := fmt.Sprintf("An alias with %v already exists\n", fmt.Sprintf("%v@%v", alias1, domain))

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can delete an user alias", func() {
		if skipUserAlias && !isCI {
			Skip("can delete an user alias")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "mailbox", "add-alias", alias1, mailboxName1)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		domain := strings.Split(mailboxName1, "@")[1]

		cli = exec.Command(cliPath, "mailbox", "delete-alias", fmt.Sprintf("%v@%v", alias1, domain))
		output, err := cli.CombinedOutput()
		Expect(err).ToNot(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Successfully deleted user alias %v\n", fmt.Sprintf("%v@%v", alias1, domain))

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can't delete an alias which doesn't exist", func() {
		if skipUserAlias && !isCI {
			Skip("can't delete an alias which doesn't exist")
		}

		cli := exec.Command(cliPath, "mailbox", "delete-alias", mailboxName1)
		output, err := cli.CombinedOutput()
		Expect(err).To(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("An alias with %v doesn't exists\n", mailboxName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

})
