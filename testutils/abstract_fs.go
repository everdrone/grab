package testutils

import "github.com/spf13/afero"

func GetOSRoot() string {
	if afero.FilePathSeparator == "\\" {
		return "C:\\"
	}
	return afero.FilePathSeparator
}

func SetupMemMapFs(root string) (afero.Fs, *afero.Afero, string) {
	fs := afero.NewMemMapFs()
	afs := &afero.Afero{Fs: fs}
	wd := root

	return fs, afs, wd
}
