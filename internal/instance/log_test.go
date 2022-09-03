package instance

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/rs/zerolog/log"
)

func TestDefaultLogger(t *testing.T) {
	tests := []struct {
		name string
		f    func()
		want string
	}{
		{
			"trace",
			func() { log.Trace().Bool("test", true).Msg("trace") },
			" TRC  trace test=true\n",
		},
		{
			"debug",
			func() { log.Debug().Bool("test", true).Msg("debug") },
			" DBG  debug test=true\n",
		},
		{
			"info",
			func() { log.Info().Bool("test", true).Msg("info") },
			" INF  info test=true\n",
		},
		{
			"warning",
			func() { log.Warn().Bool("test", true).Msg("warn") },
			" WRN  warn test=true\n",
		},
		{
			"error",
			func() { log.Err(fmt.Errorf("text")).Bool("test", true).Msg("err") },
			" ERR  err [36merror=[0m[31mtext[0m test=true\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.Buffer{}

			log.Logger = log.Output(DefaultLogger(&buf))

			tt.f()

			if got := buf.String(); got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}
