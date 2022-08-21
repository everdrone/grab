package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/everdrone/grab/internal/utils"
	tu "github.com/everdrone/grab/testutils"
)

func TestFindCmd(t *testing.T) {
	initialWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(initialWd)
	}()

	// we need to use a fake root directory for testing on windows, since filepath.Join() will not work with "\\"
	root := tu.GetOSRoot()
	utils.Fs, utils.AFS, utils.Wd = tu.SetupMemMapFs(root)

	utils.Fs.MkdirAll(filepath.Join(root, "other", "directory"), os.ModePerm)
	utils.Fs.MkdirAll(filepath.Join(root, "tmp", "test", "config", "nested"), os.ModePerm)
	utils.AFS.WriteFile(filepath.Join(root, "tmp", "test", "config", "grab.hcl"), []byte("something"), os.ModePerm)

	tests := []struct {
		Name      string
		Wd        string
		Args      []string
		Want      string
		WantErr   string
		HasErrors bool
	}{
		{
			Name: "found in current directory",
			Wd:   filepath.Join(root, "tmp", "test", "config"),
			Want: filepath.Join(root, "tmp", "test", "config", "grab.hcl") + "\n",
		},
		{
			Name: "found in parent directory",
			Wd:   filepath.Join(root, "tmp", "test", "config", "nested"),
			Want: filepath.Join(root, "tmp", "test", "config", "grab.hcl") + "\n",
		},
		{
			Name:      "not found",
			Wd:        filepath.Join(root, "other", "directory"),
			Want:      "",
			WantErr:   "could not resolve grab.hcl\n",
			HasErrors: true,
		},
		{
			Name:      "search path",
			Wd:        root,
			Args:      []string{"-p", filepath.Join(root, "tmp", "test", "config", "nested")},
			Want:      filepath.Join(root, "tmp", "test", "config", "grab.hcl") + "\n",
			WantErr:   "",
			HasErrors: false,
		},
		{
			Name:      "search path with relative path",
			Wd:        filepath.Join(root, "other", "directory"),
			Args:      []string{"-p", filepath.Join("..", "..", "tmp", "test", "config")},
			Want:      filepath.Join(root, "tmp", "test", "config", "grab.hcl") + "\n",
			WantErr:   "",
			HasErrors: false,
		},
		{
			Name:      "search path does not exist",
			Wd:        filepath.Join(root, "other", "directory"),
			Args:      []string{"-p", filepath.Join(root, "tmp", "test", "config", "nested", "not", "found")},
			Want:      "",
			WantErr:   fmt.Sprintf("path does not exist: %s\n", filepath.Join(root, "tmp", "test", "config", "nested", "not", "found")),
			HasErrors: true,
		},
	}

	args := []string{"config", "find"}

	for _, tt := range tests {
		t.Run(tt.Name, func(tc *testing.T) {
			func() {
				utils.Wd = tt.Wd
			}()

			c, got, gotErr, err := tu.ExecuteCommandErr(RootCmd, append(args, tt.Args...)...)
			if (err != nil) != tt.HasErrors {
				tc.Log(utils.Wd)
				tc.Errorf("got: %v, want: %v", err, tt.HasErrors)
			}

			if c.Name() != "find" {
				tc.Errorf("got: %s, want: 'find", c.Name())
			}
			if got != tt.Want {
				tc.Errorf("got: %s, want: %s", got, tt.Want)
			}
			if gotErr != tt.WantErr {
				tc.Errorf("got: %s, want: %s", gotErr, tt.WantErr)
			}
		})
	}
}
