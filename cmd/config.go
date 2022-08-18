package cmd

import (
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage the configuration file",
	Example: `  Check for errors in the configuration file:
    grab config check -c ../grab.hcl

  Generate the default configuration:
    grab config generate

  Print the path of the closest config file:
    grab config find`,
}

func init() {
	RootCmd.AddCommand(ConfigCmd)
}
