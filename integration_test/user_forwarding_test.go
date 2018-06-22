package integrationTest

import (
	"database/sql"
	"fmt"
	"os/exec"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("user-forwarding", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can add a user-forwarding", func() {
		cli := exec.Command(cliPath, "user", "add", userName, userPW)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "user", "add-forwarding", userName, forwardingAddress)
		output, err := cli.CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Successfully added user forwarding %v -> info@example.com\n", userName)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = '` + userName + `' AND forwarding = '` + forwardingAddress + `' 
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_maillist = 0);`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())

		Expect(exists).To(Equal(true))
	})

	It("can delete a user-forwarding", func() {
		cli := exec.Command(cliPath, "user", "add", userName, userPW)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "user", "add-forwarding", userName, forwardingAddress)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "user", "delete-forwarding", userName, forwardingAddress)
		output, err := cli.CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Successfully deleted user forwarding %v -> info@example.com\n", userName)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = '` + userName + `' AND forwarding = '` + forwardingAddress + `' 
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_maillist = 0);`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())

		Expect(exists).To(Equal(false))
	})
})
