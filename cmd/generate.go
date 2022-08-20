package cmd

import (
	"bytes"
	_ "embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/everdrone/grab/internal/utils"

	"github.com/spf13/cobra"
)

//go:embed grab.hcl.gtpl
var configTemplate string

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the default configuration file in the current directory",
	Long: `To write the file contents to stdout, use the --stdout flag.
To write the file contents to a file or into a directory, use the --output flag.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.Getwd()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		stdout, err := cmd.Flags().GetBool("stdout")
		if err != nil {
			cmd.PrintErrln(err)
			return utils.ErrSilent
		}

		outPath, err := cmd.Flags().GetString("output")
		if err != nil {
			cmd.PrintErrln(err)
			return utils.ErrSilent
		}

		if stdout && outPath != "" {
			cmd.PrintErrln("error: flag `output` is mutually exclusive with flag `stdout`")
			cmd.Usage()
			return utils.ErrSilent
		}

		tmpl, err := template.New("test").Parse(configTemplate)
		if err != nil {
			cmd.PrintErrf("could not generate default config from template: %v\n", err)
			return utils.ErrSilent
		}

		homedir, err := os.UserHomeDir()
		if err != nil {
			cmd.PrintErrf("could not get user home directory: %v\n", err)
			return utils.ErrSilent
		}

		// escape backslashes if any
		homedir = strings.Replace(homedir, "\\", "\\\\", -1)

		data := &ConfigTemplateData{
			Location: filepath.Join(homedir, "Downloads", "grab"),
		}

		buffer := new(bytes.Buffer)
		if err = tmpl.Execute(buffer, data); err != nil {
			cmd.PrintErrf("could not generate default config from template: %v\n", err)
			return utils.ErrSilent
		}

		if stdout {
			cmd.Print(buffer.String())
			return nil
		}

		if outPath == "" {
			if err = utils.AFS.WriteFile(filepath.Join(utils.Wd, "grab.hcl"), buffer.Bytes(), os.ModePerm); err != nil {
				cmd.PrintErrf("could not write config to file: %v\n", err)
				return utils.ErrSilent
			}
		} else {
			fileInfo, err := utils.Fs.Stat(outPath)
			if err != nil {
				cmd.PrintErrf("could not get file info for output path: %v\n", err)
				return utils.ErrSilent
			}

			if fileInfo.IsDir() {
				outPath = filepath.Join(outPath, "grab.hcl")
			}

			if err = utils.AFS.WriteFile(outPath, buffer.Bytes(), os.ModePerm); err != nil {
				cmd.PrintErrf("could not write config to file: %v\n", err)
				return utils.ErrSilent
			}
		}

		return nil

	},
}

type ConfigTemplateData struct {
	Location string
}

func init() {
	ConfigCmd.AddCommand(GenerateCmd)

	GenerateCmd.Flags().BoolP("stdout", "f", false, "write to stdout")
	GenerateCmd.Flags().StringP("output", "o", "", "the path to the output file")
}
