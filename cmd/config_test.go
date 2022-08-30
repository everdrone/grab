package cmd

import (
	"strings"
	"testing"

	tu "github.com/everdrone/grab/testutils"
)

func TestConfigCmd(t *testing.T) {
	t.Run("prints help message", func(tc *testing.T) {
		c, got, err := tu.ExecuteCommand(RootCmd, "config")

		if err != nil {
			tc.Fatal(err)
		}

		if c.Name() != ConfigCmd.Name() {
			tc.Fatalf("got: '%s', want: '%s'", c.Name(), ConfigCmd.Name())
		}

		if !strings.HasPrefix(got, ConfigCmd.Short) {
			tc.Errorf("got: '%s', want: '%s'", got, ConfigCmd.Short)
		}
	})
}
