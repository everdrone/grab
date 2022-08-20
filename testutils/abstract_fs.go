package testutils

import "github.com/spf13/afero"

func GetOSRoot() string {
	root := afero.FilePathSeparator
	if root == "\\" {
		root = "C:\\"
	}

	return root
}

func SetupMemMapFs(root string) (afero.Fs, *afero.Afero, string) {
	fs := afero.NewMemMapFs()
	afs := &afero.Afero{Fs: fs}
	wd := root

	return fs, afs, wd
}
