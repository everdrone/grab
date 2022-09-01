package testutils

import (
	"io/fs"
	"os"
	"strings"

	"github.com/spf13/afero"
)

func GetOSRoot() string {
	if afero.FilePathSeparator == "\\" {
		return "C:\\"
	}
	return afero.FilePathSeparator
}

type MockFs struct {
	afero.MemMapFs
}

// this function gives error when opening file_does_not_exist.txt
func (m *MockFs) Open(name string) (afero.File, error) {
	if strings.Contains(name, "restricted__r") {
		return nil, fs.ErrPermission
	}

	return m.MemMapFs.Open(name)
}

func (m *MockFs) Create(name string) (afero.File, error) {
	if strings.Contains(name, "restricted__w") {
		return nil, fs.ErrPermission
	}

	return m.MemMapFs.Create(name)
}

func (m *MockFs) MkdirAll(path string, perm fs.FileMode) error {
	if strings.Contains(path, "restricted__m") {
		return fs.ErrPermission
	}

	return m.MemMapFs.MkdirAll(path, perm)
}

type MockIoUtil struct{}

func (i *MockIoUtil) WriteFile(afs afero.Fs, filename string, data []byte, perm os.FileMode) error {
	if strings.Contains(filename, "restricted__w") {
		return fs.ErrPermission
	}

	return afero.WriteFile(afs, filename, data, perm)
}

func (i *MockIoUtil) ReadFile(afs afero.Fs, filename string) ([]byte, error) {
	if strings.Contains(filename, "restricted__r") {
		return nil, fs.ErrPermission
	}

	return afero.ReadFile(afs, filename)
}

func (i *MockIoUtil) Exists(afs afero.Fs, path string) (bool, error) {
	if strings.Contains(path, "restricted__e") {
		return false, fs.ErrPermission
	}

	return afero.Exists(afs, path)
}

func SetupMemMapFs(root string) (afero.Fs, *MockIoUtil, string) {
	fs := new(MockFs)
	wd := root

	return fs, &MockIoUtil{}, wd
}
