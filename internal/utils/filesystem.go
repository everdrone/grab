package utils

import (
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

var (
	Fs  afero.Fs     = afero.NewOsFs()
	AFS *afero.Afero = &afero.Afero{Fs: Fs}
	Wd  string       = ""
)

// sets the working directory globally
func Getwd() error {
	if Wd != "" {
		// prevents calling the function twice, actually useful when testing
		return nil
	}

	workDir, err := os.Getwd()
	if err != nil {
		return err
	}

	Wd = workDir
	return nil
}

func Abs(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(Wd, rel)
}
