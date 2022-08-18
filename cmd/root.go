package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/everdrone/grab/internal/config"
	"github.com/everdrone/grab/internal/utils"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "grab",
	// TODO: describe both short and long
	Short: "Download and scrape web pages based on a set of regex patterns",
	Example: `  Generate the starter config file:
    grab config generate

  After customizing it, check for errors:
    grab config check

  Scrape and download assets from a url:
    grab get http://example.com

  Scrape and download assets from a list of urls:
    grab get list.ini`,
	SilenceErrors: true,
	SilenceUsage:  true,
	Args:          cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vt:", cmd.VersionTemplate())
		if bi, ok := debug.ReadBuildInfo(); ok {
			fmt.Printf("\nBuild info: %+v\n", bi)
		}
		fmt.Println(config.BuildArch, config.BuildOS)
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// from: https://github.com/spf13/cobra/issues/914#issuecomment-548411337
	RootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.PrintErrf("Error: %s\n", err)
		cmd.Println(cmd.UsageString())
		return utils.ErrSilent
	})
}
