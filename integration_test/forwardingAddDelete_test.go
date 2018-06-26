package integrationTest

import (
	"database/sql"
	"fmt"
	"os/exec"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("forwarding", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can add a forwarding", func() {
		if skipUserForwarding && !isCI {
			Skip("can add a forwarding")
		}

		cli := exec.Command(cliPath, "user", "add", userName1, userPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "forwarding", "add", userName1, forwardingAddress)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := fmt.Sprintf("Successfully added forwarding %v -> info@example.com\n", userName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = '` + userName1 + `' AND forwarding = '` + forwardingAddress + `' 
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_maillist = 0);`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())

		Expect(exists).To(Equal(true))
	})

	It("can delete a forwarding", func() {
		if skipUserForwarding && !isCI {
			Skip("can delete a forwarding")
		}

		cli := exec.Command(cliPath, "user", "add", userName1, userPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "forwarding", "add", userName1, forwardingAddress)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = '` + userName1 + `' AND forwarding = '` + forwardingAddress + `' 
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_maillist = 0);`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))

		cli = exec.Command(cliPath, "forwarding", "delete", userName1, forwardingAddress)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := fmt.Sprintf("Successfully deleted forwarding %v -> info@example.com\n", userName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		query = `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = '` + userName1 + `' AND forwarding = '` + forwardingAddress + `' 
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_maillist = 0);`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())

		Expect(exists).To(Equal(false))
	})
})
