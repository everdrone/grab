package cmd

import (
	"fmt"
	"testing"

	"github.com/everdrone/grab/internal/config"
	"github.com/everdrone/grab/internal/testutils"
)

func TestVersionCmd(t *testing.T) {
	t.Run("version", func(t *testing.T) {
		cmdName := "version"

		c, got, err := testutils.ExecuteCommand(RootCmd, cmdName)
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
