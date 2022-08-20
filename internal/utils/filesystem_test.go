package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/everdrone/grab/testutils"
)

func resetWd() {
	Wd = ""
}

func TestGetwd(t *testing.T) {
	initialWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(initialWd)
	}()

	tests := []struct {
		Name    string
		Wd      string
		WantErr bool
		WantWd  string
	}{
		{
			Name:    "already set",
			Wd:      "/tmp",
			WantErr: false,
			WantWd:  "/tmp",
		},
		{
			Name:    "not set",
			Wd:      "",
			WantErr: false,
			WantWd:  initialWd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			Wd = tt.Wd

			if err := Getwd(); (err != nil) != tt.WantErr {
				t.Errorf("got: %v, want: %v", err, tt.WantErr)
			}
			if Wd != tt.WantWd {
				t.Errorf("got: %v, want: %v", Wd, tt.WantWd)
			}

			resetWd()
		})
	}
}

func TestAbs(t *testing.T) {
	initialWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(initialWd)
	}()

	root := testutils.GetOSRoot()

	tests := []struct {
		Wd    string
		Input string
		Want  string
	}{
		{
			Wd:    filepath.Join(root, "tmp"),
			Input: "config",
			Want:  filepath.Join(root, "tmp", "config"),
		},
		{
			Wd:    filepath.Join(root, "other", "directory"),
			Input: "config",
			Want:  filepath.Join(root, "other", "directory", "config"),
		},
		{
			Wd:    filepath.Join(root, "other", "directory"),
			Input: filepath.Join("..", "..", "config"),
			Want:  filepath.Join(root, "config"),
		},
		{
			Wd:    filepath.Join(root, "other", "directory"),
			Input: filepath.Join("..", "..", "..", ".."),
			Want:  filepath.Join(root),
		},
		{
			Wd:    filepath.Join(root, "other", "directory"),
			Input: filepath.Join("..", "..", "..", "..", "tmp", "dir"),
			Want:  filepath.Join(root, "tmp", "dir"),
		},
		{
			Wd:    filepath.Join(root, "other", "directory"),
			Input: filepath.Join(root, "tmp", "dir"),
			Want:  filepath.Join(root, "tmp", "dir"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Input, func(tc *testing.T) {
			Wd = tt.Wd

			got := Abs(tt.Input)
			if got != tt.Want {
				tc.Errorf("got: %v, want: %v", got, tt.Want)
			}

			if !filepath.IsAbs(got) {
				tc.Errorf("got: %v, want: %v", got, tt.Want)
			}
		})
	}

	resetWd()
}
