package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/everdrone/grab/internal/utils"
	tu "github.com/everdrone/grab/testutils"
)

func TestGenerate(t *testing.T) {
	initialWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(initialWd)
	}()

	root := tu.GetOSRoot()

	tests := []struct {
		Name      string
		Wd        string
		Args      []string
		HasErrors bool
		CheckFile string
		Want      string
	}{
		{
			Name:      "no args",
			Wd:        filepath.Join(root, "test"),
			HasErrors: false,
			CheckFile: filepath.Join(root, "test", "grab.hcl"),
		},
		{
			Name:      "stdout",
			Args:      []string{"--stdout"},
			Wd:        filepath.Join(root, "test"),
			HasErrors: false,
			Want:      ``,
		},
	}

	args := []string{"config", "generate"}

	for _, tt := range tests {
		t.Run(tt.Name, func(tc *testing.T) {
			utils.Fs, utils.AFS, utils.Wd = tu.SetupMemMapFs(root)
			utils.Fs.MkdirAll(filepath.Join(root, "test"), os.ModePerm)

			func() {
				utils.Wd = tt.Wd
			}()

			c, got, _, err := tu.ExecuteCommandErr(RootCmd, append(args, tt.Args...)...)
			if (err != nil) != tt.HasErrors {
				tc.Log(utils.Wd)
				tc.Errorf("got: %v, want: %v", err, tt.HasErrors)
			}

			if c.Name() != "generate" {
				tc.Errorf("got: %s, want: 'generate", c.Name())
			}

			if tt.CheckFile != "" {
				gotFile, err := utils.AFS.ReadFile(tt.CheckFile)
				if err != nil {
					tc.Errorf("could not read file: %v", err)
				}

				if !strings.HasPrefix(string(gotFile), "global {") {
					tc.Errorf("got: %q, want: 'global {'", string(gotFile))
				}
			} else {
				if !strings.HasPrefix(got, tt.Want) {
					tc.Errorf("got: %s, want: %s", got, tt.Want)
				}
			}
		})
	}
}
