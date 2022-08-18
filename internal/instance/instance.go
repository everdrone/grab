package instance

import (
	"github.com/everdrone/grab/internal/config"
	"github.com/spf13/cobra"
)

type FlagsState struct {
	// the path of the config file
	ConfigPath string
	// the verbosity level (0 = quiet, 1 = default, 2 = verbose, 3 = debug)
	Verbosity int
	// force overwrite of downloaded assets
	Force bool
	// stop at the first error
	Strict bool
	// skip writing to the disk
	DryRun bool
	// do not emit output
	Quiet bool
	// display the progress bar
	Progress bool
}

type Grab struct {
	// the parsed configuration
	Config *config.Config
	// the flags of the "get" command
	Flags *FlagsState
	// the caller command ("get")
	Command *cobra.Command

	// the original urls passed as arguments
	URLs []string
	// the main location to write to
	GlobalLocation string
	// the number of assets to be downloaded
	TotalAssets int64
	// a map of all the regular expressions to be used
	RegexCache config.RegexCacheMap
}

func New(cmd *cobra.Command) *Grab {
	return &Grab{
		Command: cmd,
	}
}
