package cmd

import (
	"github.com/everdrone/grab/internal/instance"
	"github.com/everdrone/grab/internal/utils"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Scrape and download assets from a URL, a file or a both",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Logger = instance.DefaultLogger(cmd.OutOrStderr())

		g := instance.New(cmd)

		g.ParseFlags()

		if diags := g.ParseConfig(); diags.HasErrors() {
			for _, diag := range diags.Errs() {
				log.Err(diag).Msg("config error")
			}
			return utils.ErrSilent
		}

		if diags := g.ParseURLs(args); diags.HasErrors() {
			for _, diag := range diags.Errs() {
				log.Err(diag).Msg("argument error")
			}
			return utils.ErrSilent
		}

		g.BuildSiteCache()

		if diags := g.BuildAssetCache(); diags.HasErrors() {
			for _, diag := range diags.Errs() {
				log.Err(diag).Msg("runtime error")
			}
			return utils.ErrSilent
		}

		if err := g.Download(); err != nil {
			return utils.ErrSilent
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
