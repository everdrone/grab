package cmd

import (
	"fmt"
	"testing"

	"github.com/everdrone/grab/internal/config"
	tu "github.com/everdrone/grab/testutils"
)

func TestVersionCmd(t *testing.T) {
	t.Run("version", func(t *testing.T) {
		cmdName := "version"

		c, got, err := tu.ExecuteCommand(RootCmd, cmdName)
		if err != nil {
			t.Fatal(err)
		}

		want := fmt.Sprintf("%s v%s %s/%s (%s)\n",
			"grab",
			config.Version,
			"unknown",
			"unknown",
			"unknown",
		)
		if c.Name() != cmdName {
			t.Fatalf("got: '%s', want: '%s'", c.Name(), cmdName)
		}
		if got != want {
			t.Errorf("got: %s, want: %s", got, want)
		}
	})
}
