package instance

import (
	"testing"

	"github.com/spf13/cobra"
)

func createMockGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Just a mock",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().BoolP("force", "f", false, "overwrite existing files")
	cmd.Flags().StringP("config", "c", "", "the path of the config file to use")

	cmd.Flags().BoolP("strict", "s", false, "fail on errors")
	cmd.Flags().BoolP("dry-run", "n", false, "do not write on disk")

	cmd.Flags().BoolP("progress", "p", false, "show progress bars")
	cmd.Flags().BoolP("quiet", "q", false, "do not emit any output")
	cmd.Flags().CountP("verbose", "v", "verbosity level")

	return cmd
}

func TestNew(t *testing.T) {
	t.Run("new returns instance", func(tc *testing.T) {
		mock := createMockGetCmd()

		got := New(mock)
		if got == nil {
			tc.Errorf("got: %v, want: %v", got, nil)
		}

		if got.Command != mock {
			tc.Errorf("got: %v, want: %v", got.Command, mock)
		}
	})
}

func TestParseFlags(t *testing.T) {
	t.Run("parse quiet", func(tc *testing.T) {
		mock := createMockGetCmd()
		g := New(mock)

		mock.SetArgs([]string{"http://example.com", "-vvv", "-q"})
		mock.ExecuteC()

		g.ParseFlags()

		if g.Flags.Quiet != true {
			tc.Errorf("got: %v, want: %v", g.Flags.Quiet, true)
		}
		if g.Flags.Verbosity != 0 {
			tc.Errorf("got: %v, want: %v", g.Flags.Verbosity, 0)
		}
	})

	t.Run("parse verbosity", func(tc *testing.T) {
		mock := createMockGetCmd()
		g := New(mock)

		mock.SetArgs([]string{"http://example.com", "-vv"})
		mock.ExecuteC()

		g.ParseFlags()

		if g.Flags.Verbosity != 3 {
			tc.Errorf("got: %v, want: %v", g.Flags.Verbosity, 3)
		}
	})
}

// func TestParseConfig(t *testing.T) {
// 	t.Fail()
// }

// func TestParseURLs(t *testing.T) {
// 	t.Fail()
// }

// func TestBuildSiteCache(t *testing.T) {
// 	t.Fail()
// }
