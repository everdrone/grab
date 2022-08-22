package cmd

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/everdrone/grab/internal/config"
	tu "github.com/everdrone/grab/testutils"
)

func TestVersionCmd(t *testing.T) {
	cmdName := "version"

	config.CommitHash = "abcdef0123456789"
	config.BuildOS = runtime.GOOS
	config.BuildArch = runtime.GOARCH

	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		want    string
	}{
		{
			name: "no updates",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{"tag_name": "v` + config.Version + `"}`))
			},
			want: "grab v" + config.Version + " " + config.BuildOS + "/" + config.BuildArch + " (" + config.CommitHash[:7] + ")\n",
		},
		{
			name: "newer version",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{"tag_name": "v987.654.321"}`))
			},
			want: "grab v" + config.Version + " " + config.BuildOS + "/" + config.BuildArch + " (" + config.CommitHash[:7] + ")\n\n" +
				"New version available: v987.654.321\n" +
				"Download it at https://github.com/everdrone/grab/releases/latest\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(tc *testing.T) {
			// start the test server
			ts := httptest.NewServer(http.HandlerFunc(tt.handler))

			config.LatestReleaseURL = ts.URL

			c, got, err := tu.ExecuteCommand(RootCmd, cmdName)
			if err != nil {
				tc.Fatal(err)
			}

			if c.Name() != cmdName {
				tc.Fatalf("got: '%s', want: '%s'", c.Name(), cmdName)
			}

			if got != tt.want {
				tc.Errorf("got: '%s', want: '%s'", got, tt.want)
			}
		})
	}
}
