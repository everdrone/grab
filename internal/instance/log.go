package instance

import (
	"github.com/everdrone/grab/internal/utils"

	"github.com/hashicorp/hcl/v2"
)

func (s *Grab) Log(level int, diags hcl.Diagnostics) {
	for _, diag := range diags {
		switch diag.Severity {
		case utils.DiagError:
			fallthrough
		case utils.DiagWarning:
			utils.PrintDiag(s.Command.ErrOrStderr(), diag)
		default:
			if s.Flags.Verbosity >= level {
				utils.PrintDiag(s.Command.OutOrStdout(), diag)
			}
		}
	}
}

/*
╷ Error: Something bad happened
│   Description of the error, here we can be more detailed
╵   something.hcl:2-3
╷ Error: Something bad happened
╵   Description of the error, here we can be more detailed
*/
