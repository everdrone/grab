package cmd

import (
	"github.com/everdrone/grab/internal/config"
	"github.com/everdrone/grab/internal/instance"
	"github.com/everdrone/grab/internal/update"
	"github.com/everdrone/grab/internal/utils"
	"github.com/fatih/color"

	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Scrape and download assets from a URL, a file or a both",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		updateMessageChan := make(chan string)
		go func() {
			newVersion, _ := update.CheckForUpdates()
			updateMessageChan <- newVersion
		}()

		g := instance.New(cmd)

		g.ParseFlags()

		if diags := g.ParseConfig(); diags.HasErrors() {
			g.Log(0, *diags)
			return utils.ErrSilent
		}

		if diags := g.ParseURLs(args); diags.HasErrors() {
			g.Log(0, *diags)
			return utils.ErrSilent
		}

		g.BuildSiteCache()

		if diags := g.BuildAssetCache(); diags.HasErrors() {
			g.Log(0, *diags)
			return utils.ErrSilent
		}

		if diags := g.Download(); diags.HasErrors() {
			g.Log(0, *diags)
			return utils.ErrSilent
		}

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

		return nil
	},
}

func init() {
	RootCmd.AddCommand(GetCmd)

	GetCmd.Flags().BoolP("force", "f", false, "overwrite existing files")
	GetCmd.Flags().StringP("config", "c", "", "the path of the config file to use")

	GetCmd.Flags().BoolP("strict", "s", false, "fail on errors")
	GetCmd.Flags().BoolP("dry-run", "n", false, "do not write on disk")

	GetCmd.Flags().BoolP("progress", "p", false, "show progress bars")
	GetCmd.Flags().BoolP("quiet", "q", false, "do not emit any output")
	GetCmd.Flags().CountP("verbose", "v", "verbosity level")
}
