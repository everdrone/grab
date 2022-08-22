package cmd

import (
	"github.com/everdrone/grab/internal/config"
	"github.com/everdrone/grab/internal/update"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number and check for updates",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("%s v%s %s/%s (%s)\n",
			cmd.Root().Name(),
			config.Version,
			config.BuildOS,
			config.BuildArch,
			config.CommitHash[:7])

		newVersion, _ := update.CheckForUpdates()
		if newVersion != "" {
			cmd.Printf("\nNew version available %s → %s\n", config.Version, newVersion)
			cmd.Printf("https://github.com/everdrone/grab/releases/latest\n")
		}
	},
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}
