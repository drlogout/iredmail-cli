package integrationTest

import (
	"fmt"
	"os/exec"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("alias info", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can show aliases info", func() {
		if skipAliasInfo && !isCI {
			Skip("can show aliases info")
		}

		aliases := []string{
			alias1,
			alias2,
			alias3,
		}

		for _, alias := range aliases {
			cli := exec.Command(cliPath, "alias", "add", alias)
			output, err := cli.CombinedOutput()
			if err != nil {
				Fail(string(output))
			}
		}

		cli := exec.Command(cliPath, "alias", "add-forwarding", alias1, aliasForwarding1)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "alias", "add-forwarding", alias1, aliasForwarding2)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "alias", "add-forwarding", alias2, aliasForwarding2)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "alias", "info", alias1)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := loadGolden("can_show_alias_info")

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})
})
