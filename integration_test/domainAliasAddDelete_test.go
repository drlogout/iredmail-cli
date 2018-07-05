package integrationTest

import (
	"database/sql"
	"fmt"
	"os/exec"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("domain add/delete-alias", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can add an alias domain", func() {
		if skipDomainAliasAddDelete && !isCI {
			Skip("can add an alias domain")
		}

		cli := exec.Command(cliPath, "domain", "add", domain1)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "add-alias", aliasDomain1, domain1)
		output, err := cli.CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Successfully added alias domain %s ➞ %s\n", aliasDomain1, domain1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT alias_domain FROM alias_domain
		WHERE alias_domain = '` + aliasDomain1 + `' AND target_domain = '` + domain1 + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))
	})

	It("can't add an existing alias domain", func() {
		if skipDomainAliasAddDelete && !isCI {
			Skip("can't add an existing alias domain")
		}

		cli := exec.Command(cliPath, "domain", "add", domain1)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "add-alias", aliasDomain1, domain1)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "add-alias", aliasDomain1, domain1)
		output, err := cli.CombinedOutput()
		Expect(err).To(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Alias domain %s ➞ %s already exists\n", aliasDomain1, domain1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can't add an alias domain if domain doesn't exist", func() {
		if skipDomainAliasAddDelete && !isCI {
			Skip("can't add an alias domain if domain doesn't exist")
		}

		cli := exec.Command(cliPath, "domain", "add-alias", aliasDomain1, domain1)
		output, err := cli.CombinedOutput()
		Expect(err).To(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Domain %s doesn't exists\n", domain1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can delete an alias domain", func() {
		if skipDomainAliasAddDelete && !isCI {
			Skip("can delete an alias domain")
		}

		cli := exec.Command(cliPath, "domain", "add", domain1)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "add-alias", aliasDomain1, domain1)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "delete-alias", aliasDomain1)
		output, err := cli.CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Successfully deleted alias domain %s\n", aliasDomain1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT alias_domain FROM alias_domain
		WHERE alias_domain = '` + aliasDomain1 + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(false))
	})

	It("can't delete a non existing alias domain", func() {
		if skipDomainAliasAddDelete && !isCI {
			Skip("can't delete a non existing alias domain")
		}

		cli := exec.Command(cliPath, "domain", "add", domain1)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "delete-alias", aliasDomain1)
		output, err := cli.CombinedOutput()
		Expect(err).To(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Alias domain %s doesn't exist\n", aliasDomain1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})
})
