package utils

import (
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

var (
	Fs afero.Fs       = afero.NewOsFs()
	Io AbstractIoUtil = &IoUtil{}
	Wd string         = ""
)

type AbstractIoUtil interface {
	WriteFile(fs afero.Fs, filename string, data []byte, perm os.FileMode) error
	ReadFile(fs afero.Fs, filename string) ([]byte, error)
	Exists(fs afero.Fs, path string) (bool, error)
}

type IoUtil struct{}

func (i *IoUtil) WriteFile(fs afero.Fs, filename string, data []byte, perm os.FileMode) error {
	return afero.WriteFile(fs, filename, data, perm)
}

func (i *IoUtil) ReadFile(fs afero.Fs, filename string) ([]byte, error) {
	return afero.ReadFile(fs, filename)
}

func (i *IoUtil) Exists(fs afero.Fs, path string) (bool, error) {
	return afero.Exists(fs, path)
}

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
