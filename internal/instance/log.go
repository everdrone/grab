package instance

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func DefaultLogger(w io.Writer) zerolog.Logger {
	return log.Output(zerolog.ConsoleWriter{
		Out:             w,
		TimeFormat:      time.RFC3339Nano,
		FormatTimestamp: func(i interface{}) string { return "" },
		FormatLevel: func(i interface{}) string {
			var l string
			c := color.New(color.FgBlack)
			if ll, ok := i.(string); ok {
				switch ll {
				case zerolog.LevelTraceValue:
					c.Add(color.BgMagenta)
					l = c.Sprint(" TRC ")
				case zerolog.LevelDebugValue:
					c.Add(color.BgBlue)
					l = c.Sprint(" DBG ")
				case zerolog.LevelInfoValue:
					c.Add(color.BgGreen)
					l = c.Sprint(" INF ")
				case zerolog.LevelWarnValue:
					c.Add(color.BgYellow)
					l = c.Sprint(" WRN ")
				case zerolog.LevelErrorValue:
					c.Add(color.BgRed)
					l = c.Sprint(" ERR ")
				case zerolog.LevelFatalValue:
					c.Add(color.BgHiRed)
					l = c.Sprint(" FTL ")
				case zerolog.LevelPanicValue:
					c.Add(color.BgHiRed).Add(color.Bold)
					l = c.Sprint(" PNC ")
				default:
					c.Add(color.FgWhite).Add(color.Bold)
					l = c.Sprint(" ??? ")
				}
			} else {
				if i == nil {
					c.Add(color.FgWhite).Add(color.Bold)
					l = c.Sprint(" ??? ")
				} else {
					l = " " + strings.ToUpper(fmt.Sprintf("%s", i))[0:5] + " "
				}
			}
			return l
		},
		FormatFieldName: func(i interface{}) string {
			return color.New(color.FgHiWhite).Sprint(fmt.Sprintf("%s=", i))
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		},
	})
}

/*
╷ Error: Something bad happened
│   Description of the error, here we can be more detailed
╵   something.hcl:2-3
╷ Error: Something bad happened
╵   Description of the error, here we can be more detailed
*/
