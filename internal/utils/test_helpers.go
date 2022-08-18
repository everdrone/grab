package utils

import (
	"bytes"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func String(v string) *string {
	return &v
}

func Int(v int) *int {
	return &v
}

func Bool(v bool) *bool {
	return &v
}

func ExecuteCommand(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	sout := new(bytes.Buffer)
	root.SetOut(sout)
	root.SetErr(sout)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, sout.String(), err
}

func ExecuteCommandErr(root *cobra.Command, args ...string) (c *cobra.Command, output string, errOutput string, err error) {
	sout := new(bytes.Buffer)
	serr := new(bytes.Buffer)

	root.SetOut(sout)
	root.SetErr(serr)

	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, sout.String(), serr.String(), err
}

func AssertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

func GetOSRoot() string {
	root := afero.FilePathSeparator
	if root == "\\" {
		root = "C:\\"
	}

	return root
}

func SetupMemMapFs(root string) {
	Fs = afero.NewMemMapFs()
	AFS = &afero.Afero{Fs: Fs}
	Wd = root
}
