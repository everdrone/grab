package cmd

import (
	"path/filepath"

	"github.com/everdrone/grab/internal/config"
	"github.com/everdrone/grab/internal/utils"

	"github.com/spf13/cobra"
)

var FindCmd = &cobra.Command{
	Use:   "find",
	Short: "Print the path of the closest configuration file",
	Long: `By default, this program will search in all parent directories
of the current directory for a configuration file.
To specify a path use the --path flag.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.Getwd()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		startPath, _ := cmd.Flags().GetString("path")

		if startPath == "" {
			startPath = utils.Wd
		} else {
			if !filepath.IsAbs(startPath) {
				startPath = utils.Abs(startPath)
			}

			if exists, err := utils.Io.Exists(utils.Fs, startPath); err != nil || !exists {
				cmd.PrintErrf("path does not exist: %s\n", startPath)
				return utils.ErrSilent
			}
		}

		resolved, err := config.Resolve("grab.hcl", startPath)
		if err != nil {
			cmd.PrintErrln(err)
			return utils.ErrSilent
		}

		cmd.Println(resolved)
		return nil
	},
}

func init() {
	ConfigCmd.AddCommand(FindCmd)

	FindCmd.Flags().StringP("path", "p", "", "the path to start the search from")
}
