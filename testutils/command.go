package testutils

import (
	"bytes"

	"github.com/spf13/cobra"
)

func ExecuteCommand(rootCommand *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	stdOut := new(bytes.Buffer)
	rootCommand.SetOut(stdOut)
	rootCommand.SetErr(stdOut)
	rootCommand.SetArgs(args)

	c, err = rootCommand.ExecuteC()

	return c, stdOut.String(), err
}

func ExecuteCommandErr(rootCommand *cobra.Command, args ...string) (c *cobra.Command, output string, errOutput string, err error) {
	stdOut := new(bytes.Buffer)
	stdErr := new(bytes.Buffer)

	rootCommand.SetOut(stdOut)
	rootCommand.SetErr(stdErr)

	rootCommand.SetArgs(args)

	c, err = rootCommand.ExecuteC()

	return c, stdOut.String(), stdErr.String(), err
}
