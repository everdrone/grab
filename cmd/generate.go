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
	Long:  `To write the file contents to stdout, use the --stdout flag`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.Getwd()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		stdout, _ := cmd.Flags().GetBool("stdout")

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
		tmpl := template.Must(template.New("test").Parse(configTemplate))
		if err = tmpl.Execute(buffer, data); err != nil {
			cmd.PrintErrf("could not generate default config from template: %v\n", err)
			return utils.ErrSilent
		}

		if stdout {
			cmd.Print(buffer.String())
			return nil
		}

		outPath := filepath.Join(utils.Wd, "grab.hcl")
		if err = utils.Io.WriteFile(utils.Fs, outPath, buffer.Bytes(), os.ModePerm); err != nil {
			cmd.PrintErrf("could not write config to file: %v\n", err)
			return utils.ErrSilent
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
}
