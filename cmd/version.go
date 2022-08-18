package cmd

import (
	"github.com/everdrone/grab/internal/config"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number and exit",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("%s v%s %s/%s (%s)\n",
			cmd.Root().Name(),
			config.Version,
			config.BuildOS,
			config.BuildArch,
			config.CommitHash[:7])
	},
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}
