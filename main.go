// to build on windows (powershell)
// $VHASH = git rev-parse HEAD; go build -ldflags="-s -w -X grab/internal/config.CommitHash=${VHASH}"

// to build on linux/macos
// VHASH=$(git rev-parse HEAD) go build -ldflags="-s -w -X grab/internal/config.CommitHash=${VHASH}"

package main

import (
	"fmt"
	"os"

	"github.com/everdrone/grab/cmd"
	"github.com/everdrone/grab/internal/utils"
)

func main() {
	// from: https://github.com/spf13/cobra/issues/914#issuecomment-548411337
	if err := cmd.RootCmd.Execute(); err != nil {
		// if we have ErrSilent, we don't want to print the error
		if err != utils.ErrSilent {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
