package integrationTest

import (
	"database/sql"
	"fmt"
	"os/exec"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	mailboxName       = "post@domain.com"
	mailboxPW         = "alskdlqkdjalskd"
	forwardingAddress = "info@example.com"
)

var _ = Describe("Mailbox", func() {
	It("can add a mailbox", func() {
		cli := exec.Command(cliPath, "mailbox", "add", mailboxName, mailboxPW)
		output, err := cli.CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Successfully added mailbox %v (quota: 2048 KB)\n", mailboxName)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT * FROM mailbox
		WHERE username = '` + mailboxName + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))

		query = `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = '` + mailboxName + `' AND forwarding = '` + mailboxName + `' 
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_maillist = 0);`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))
	})

	It("can delete a mailbox", func() {
		cli := exec.Command(cliPath, "mailbox", "delete", "--force", mailboxName)
		output, err := cli.CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Successfully deleted mailbox %v\n", mailboxName)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT * FROM mailbox
		WHERE username = '` + mailboxName + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(false))

		query = `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = '` + mailboxName + `' AND forwarding = '` + mailboxName + `' 
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_maillist = 0);`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(false))
	})

	It("can add an mailbox-forwarding", func() {
		cli := exec.Command(cliPath, "mailbox", "add", mailboxName, mailboxPW)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "mailbox", "add-forwarding", mailboxName, forwardingAddress)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = '` + mailboxName + `' AND forwarding = '` + forwardingAddress + `' 
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_maillist = 0);`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())

		Expect(exists).To(Equal(true))
	})
})
