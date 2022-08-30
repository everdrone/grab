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
			want: "grab v" + config.Version + " " + config.BuildOS + "/" + config.BuildArch + " (" + config.CommitHash[:7] + ")\n\n\n" +
				"A new release of grab is available: " + config.Version + " → 987.654.321\n" +
				"https://github.com/everdrone/grab/releases/latest\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(tc *testing.T) {
			// start the test server
			ts := httptest.NewServer(http.HandlerFunc(tt.handler))

			config.LatestReleaseURL = ts.URL

			c, got, err := tu.ExecuteCommand(RootCmd, "version")
			if err != nil {
				tc.Fatal(err)
			}

			if c.Name() != VersionCmd.Name() {
				tc.Fatalf("got: %s, want: %s", c.Name(), VersionCmd.Name())
			}

			if got != tt.want {
				tc.Errorf("got: %q, want: %q", got, tt.want)
			}
		})
	}
}
