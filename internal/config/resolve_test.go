package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/everdrone/grab/internal/utils"
	tu "github.com/everdrone/grab/testutils"
)

func TestResolve(t *testing.T) {
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
		Name     string
		Filename string
		Wd       string
		Want     string
		Err      string
	}{
		{
			Name:     "found in current directory",
			Wd:       filepath.Join(root, "tmp", "test", "config"),
			Filename: "grab.hcl",
			Want:     filepath.Join(root, "tmp", "test", "config", "grab.hcl"),
			Err:      "",
		},
		{
			Name:     "found in parent directory",
			Wd:       filepath.Join(root, "tmp", "test", "config", "nested"),
			Filename: "grab.hcl",
			Want:     filepath.Join(root, "tmp", "test", "config", "grab.hcl"),
			Err:      "",
		},
		{
			Name:     "not found",
			Wd:       filepath.Join(root, "tmp", "test", "config", "nested"),
			Filename: "notFound.hcl",
			Want:     "",
			Err:      "could not resolve notFound.hcl",
		},
		{
			Name:     "invalid filename",
			Wd:       filepath.Join(root, "tmp", "test", "config", "nested"),
			Filename: "\000x",
			Want:     "",
			Err:      "could not resolve \000x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if !filepath.IsAbs(tt.Wd) {
				t.Fail()
			}

			func() {
				utils.Wd = tt.Wd
			}()

			got, err := Resolve(tt.Filename, tt.Wd)

			if err != nil && tt.Err == "" {
				// we want this function not to throw but it did
				t.Errorf("got error %q", err)
			}
			if tt.Err != "" {
				if !strings.HasPrefix(fmt.Sprint(err), tt.Err) {
					t.Errorf("got: %q, want: %q", err, tt.Err)
				}
			}
			if tt.Want != "" {
				if got != tt.Want {
					t.Errorf("got: %s, want %s", got, tt.Want)
				}
			}
		})
	}
}
