package integrationTest

import (
	"fmt"
	"os/exec"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("alias list", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can list aliases", func() {
		if skipAliasList && !isCI {
			Skip("can list aliases")
		}

		aliases := []string{
			alias1,
			alias2,
			alias3,
		}

		for _, alias := range aliases {
			cli := exec.Command(cliPath, "alias", "add", alias)
			err := cli.Run()
			Expect(err).NotTo(HaveOccurred())
		}

		cli := exec.Command(cliPath, "alias", "add-forwarding", alias1, aliasForwarding1)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "alias", "add-forwarding", alias1, aliasForwarding2)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "alias", "add-forwarding", alias2, aliasForwarding3)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "alias", "list")
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := loadGolden("can_list_aliases")

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can list aliases and filter result", func() {
		if skipAliasList && !isCI {
			Skip("can list aliases and filter result")
		}

		aliases := []string{
			alias1,
			alias2,
			alias3,
		}

		for _, alias := range aliases {
			cli := exec.Command(cliPath, "alias", "add", alias)
			err := cli.Run()
			Expect(err).NotTo(HaveOccurred())
		}

		cli := exec.Command(cliPath, "alias", "add-forwarding", alias1, aliasForwarding1)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "alias", "add-forwarding", alias1, aliasForwarding2)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "alias", "add-forwarding", alias2, aliasForwarding3)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "alias", "list", "-f", "domain.com")
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := loadGolden("can_list_aliases_and_filter_result")

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})
})
