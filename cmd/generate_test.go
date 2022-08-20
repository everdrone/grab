package cmd

import (
	"os"
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
	utils.Fs, utils.AFS, utils.Wd = tu.SetupMemMapFs(root)

	tests := []struct {
		Name      string
		Wd        string
		Args      []string
		HasErrors bool
		CheckFile string
		Want      string
		WantErr   string
	}{
		{},
	}

	args := []string{"config", "generate"}

	for _, tt := range tests {
		t.Run(tt.Name, func(tc *testing.T) {
			func() {
				utils.Wd = tt.Wd
			}()

			c, got, gotErr, err := tu.ExecuteCommandErr(RootCmd, append(args, tt.Args...)...)
			if (err != nil) != tt.HasErrors {
				t.Log(utils.Wd)
				t.Errorf("got: %v, want: %v", err, tt.HasErrors)
			}

			if c.Name() != "generate" {
				t.Errorf("got: %s, want: 'generate", c.Name())
			}

			if got != tt.Want {
				t.Errorf("got: %s, want: %s", got, tt.Want)
			}
			if gotErr != tt.WantErr {
				t.Errorf("got: %s, want: %s", gotErr, tt.WantErr)
			}
		})
	}
}
