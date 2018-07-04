package integrationTest

import (
	"database/sql"
	"fmt"
	"os/exec"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("alias add/delete-forwarding", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can add an alias with forwardings", func() {
		if skipAliasForwardingAddDelete && !isCI {
			Skip("can add an alias forwardings")
		}

		cli := exec.Command(cliPath, "alias", "add", alias1)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "alias", "add-forwarding", alias1, aliasForwarding1)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := fmt.Sprintf("Successfully added alias forwarding %s ➞ %s\n", alias1, aliasForwarding1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT address FROM alias
		WHERE address = '` + alias1 + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))

		query = `SELECT exists
		(SELECT address FROM forwardings
		WHERE address = ? AND forwarding = ? AND is_list = 1 AND is_forwarding = 0 AND is_alias = 0);`

		err = db.QueryRow(query, alias1, aliasForwarding1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))
	})

	It("can't add an existing alias forwarding", func() {
		if skipAliasForwardingAddDelete && !isCI {
			Skip("can't add an existing alias forwarding")
		}

		cli := exec.Command(cliPath, "alias", "add", alias1)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "alias", "add-forwarding", alias1, aliasForwarding1)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "alias", "add-forwarding", alias1, aliasForwarding1)
		output, err := cli.CombinedOutput()
		Expect(err).To(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Alias forwarding %s ➞ %s already exists\n", alias1, aliasForwarding1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can delete an alias forwarding", func() {
		if skipAliasForwardingAddDelete && !isCI {
			Skip("can delete an alias forwarding")
		}

		cli := exec.Command(cliPath, "alias", "add", alias1)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "alias", "add-forwarding", alias1, aliasForwarding1)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "alias", "add-forwarding", alias1, aliasForwarding2)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "alias", "delete-forwarding", alias1, aliasForwarding2)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := fmt.Sprintf("Successfully deleted alias forwarding %s ➞ %s\n", alias1, aliasForwarding2)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT address FROM alias
		WHERE address = '` + alias1 + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))

		query = `SELECT exists
		(SELECT address FROM forwardings
		WHERE address = ? AND forwarding = ? AND is_list = 1 AND is_forwarding = 0 AND is_alias = 0);`

		err = db.QueryRow(query, alias1, aliasForwarding1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))

		query = `SELECT exists
		(SELECT address FROM forwardings
		WHERE address = ? AND forwarding = ? AND is_list = 1 AND is_forwarding = 0 AND is_alias = 0);`

		err = db.QueryRow(query, alias1, aliasForwarding2).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(false))
	})

	It("can delete an non existing alias forwarding", func() {
		if skipAliasForwardingAddDelete && !isCI {
			Skip("can delete an non existing alias forwarding")
		}

		cli := exec.Command(cliPath, "alias", "add", alias1)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "alias", "delete-forwarding", alias1, aliasForwarding1)
		output, err := cli.CombinedOutput()
		Expect(err).To(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("An alias forwarding %s ➞ %s doesn't exists\n", alias1, aliasForwarding1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})
})
