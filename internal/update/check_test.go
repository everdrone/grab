package update

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/everdrone/grab/internal/config"
)

func TestCheckForUpdates(t *testing.T) {
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		want    string
		wantErr bool
	}{
		{
			name: "no updates",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{"tag_name": "v` + config.Version + `"}`))
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "invalid semver",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{"tag_name": "newVersion"}`))
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{"something": "else"}`))
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid json",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{"something`))
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "empty tag name",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{"tag_name": ""}`))
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "network error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "request times out",
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(time.Millisecond * 1500)

				w.Write([]byte(`{"tag_name": "v` + config.Version + `"}`))
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "update available",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{"tag_name": "v987.654.321"}`))
			},
			want:    "v987.654.321",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(tc *testing.T) {
			// start the test server
			ts := httptest.NewServer(http.HandlerFunc(tt.handler))

			config.LatestReleaseURL = ts.URL

			got, err := CheckForUpdates()
			if (err != nil) != tt.wantErr {
				tc.Errorf("got error: '%v', want error: '%v'", err, tt.wantErr)
			}
			if got != tt.want {
				tc.Errorf("got: '%s', want: '%s'", got, tt.want)
			}
		})
	}
}
