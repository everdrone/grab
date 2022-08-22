package cmd

import (
	"strings"
	"testing"

	tu "github.com/everdrone/grab/testutils"
)

func TestRootCmd(t *testing.T) {
	t.Run("prints help message", func(tc *testing.T) {
		c, got, err := tu.ExecuteCommand(RootCmd, "")

		if err != nil {
			tc.Fatal(err)
		}

		if c.Name() != RootCmd.Name() {
			tc.Fatalf("got: '%s', want: '%s'", c.Name(), RootCmd.Name())
		}

		if !strings.HasPrefix(got, RootCmd.Short) {
			tc.Errorf("got: '%s', want: '%s'", got, RootCmd.Short)
		}
	})
}
