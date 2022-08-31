package testutils

import (
	"io/fs"
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
	if strings.HasSuffix(name, "file_does_not_exist.txt") {
		return nil, fs.ErrNotExist
	}

	if strings.HasSuffix(name, "file_not_readable.txt") {
		return nil, fs.ErrPermission
	}

	return m.MemMapFs.Open(name)
}

func (m *MockFs) Create(name string) (afero.File, error) {
	if strings.HasSuffix(name, "file_not_writable.txt") {
		return nil, fs.ErrPermission
	}

	return m.MemMapFs.Create(name)
}

type MockAfero struct {
	afero.Afero
}

func (m *MockAfero) MkdirAll(path string, perm fs.FileMode) error {
	if strings.HasSuffix(path, "dir_not_writable") {
		return fs.ErrPermission
	}

	return m.Afero.MkdirAll(path, perm)
}

func SetupMemMapFs(root string) (afero.Fs, *afero.Afero, string) {
	fs := new(MockFs)
	afs := MockAfero{afero.Afero{Fs: fs}}
	wd := root

	return fs, &afs.Afero, wd
}
