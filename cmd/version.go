package cmd

import (
	"github.com/everdrone/grab/internal/config"
	"github.com/everdrone/grab/internal/update"
	"github.com/fatih/color"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number and check for updates",
	Run: func(cmd *cobra.Command, args []string) {
		updateMessageChan := make(chan string)
		go func() {
			newVersion, _ := update.CheckForUpdates()
			updateMessageChan <- newVersion
		}()

		cmd.Printf("%s v%s %s/%s (%s)\n",
			cmd.Root().Name(),
			config.Version,
			config.BuildOS,
			config.BuildArch,
			config.CommitHash[:7])

		newVersion := <-updateMessageChan
		if newVersion != "" {
			// TODO: take in account possible package managers
			// if for example we installed with homebrew, we should display a different message
			cmd.Printf("\n\n%s %s → %s\n",
				color.New(color.FgMagenta).Sprintf("A new release of %s is available:", config.Name),
				config.Version,
				// color.New(color.FgHiBlack).Sprint(config.Version),
				color.New(color.FgCyan).Sprint(newVersion),
			)
			cmd.Printf("%s\n\n", "https://github.com/everdrone/grab/releases/latest")
		}
	},
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}
