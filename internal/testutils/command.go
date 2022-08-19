package testutils

import (
	"bytes"

	"github.com/spf13/cobra"
)

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
