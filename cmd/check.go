/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/everdrone/grab/internal/config"
	"github.com/everdrone/grab/internal/utils"

	"github.com/hashicorp/hcl/v2"
	"github.com/spf13/cobra"
)

var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate the configuration file",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.Getwd()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		quiet, _ := cmd.Flags().GetBool("quiet")
		configPath, _ := cmd.Flags().GetString("config")

		if configPath == "" {
			configPath = utils.Wd

			resolved, err := config.Resolve("grab.hcl", configPath)
			if err != nil {

				utils.PrintDiag(cmd.ErrOrStderr(), &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "could not resolve config file",
					Detail:   err.Error(),
				})
				return utils.ErrSilent
			}

			configPath = resolved
		}

		// parse config
		fc, err := utils.Io.ReadFile(utils.Fs, configPath)
		if err != nil {
			utils.PrintDiag(cmd.ErrOrStderr(), &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "could not read config file",
				Detail:   err.Error(),
			})
			return utils.ErrSilent
		}

		_, _, _, diags := config.Parse(fc, configPath)
		if diags.HasErrors() {
			for _, diag := range diags {
				utils.PrintDiag(cmd.ErrOrStderr(), diag)
			}
			return utils.ErrSilent
		}

		if !quiet {
			cmd.Println("ok")
		}
		return nil
	},
}

func init() {
	ConfigCmd.AddCommand(CheckCmd)

	CheckCmd.Flags().BoolP("quiet", "q", false, "do not emit any output")
	CheckCmd.Flags().StringP("config", "c", "", "the path of the config file to use")
}
