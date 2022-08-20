package instance

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/everdrone/grab/internal/config"
	"github.com/everdrone/grab/internal/utils"
	tu "github.com/everdrone/grab/testutils"
)

func TestDownload(t *testing.T) {
	initialWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(initialWd)
	}()
	root := tu.GetOSRoot()

	global := filepath.Join(root, "global")
	escapedGlobal := tu.EscapeHCLString(global)

	e := tu.CreateMockServer()
	ts := httptest.NewUnstartedServer(e)
	ts.Listener.Close()
	ts.Listener = e.Listener
	ts.Start()

	defer ts.Close()

	tests := []struct {
		Name   string
		Flags  *FlagsState
		Config string
		// filepath -> contents
		Want    map[string]string
		WantErr bool
	}{
		{
			Name:  "dry run",
			Flags: &FlagsState{DryRun: true},
			Config: `
global {
	location = "` + escapedGlobal + `"
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"
	asset "image" {
		pattern = "<img src=\"([^\"]+/img/[^\"]+)"
		capture = 1
		find_all = true
	}
}`,
			WantErr: false,
		},
		{
			Name:  "simple",
			Flags: &FlagsState{},
			Config: `
global {
	location = "` + escapedGlobal + `"
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"
	asset "image" {
		pattern = "<img src=\"([^\"]+/img/[^\"]+)"
		capture = 1
		find_all = true
	}
}`,
			Want: map[string]string{
				filepath.Join(global, "example", "a.jpg"): "imagea",
				filepath.Join(global, "example", "b.jpg"): "imageb",
				filepath.Join(global, "example", "c.jpg"): "imagec",
			},
			WantErr: false,
		},
		{
			Name:  "broken urls",
			Flags: &FlagsState{Strict: true},
			Config: `
global {
	location = "` + escapedGlobal + `"
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"
	asset "broken" {
		pattern = "<a href=\"([^\"]*/broken/[^\"]+)"
		capture = 1
		find_all = true
	}
}`,
			Want:    map[string]string{},
			WantErr: true,
		},
		{
			Name:  "checks headers",
			Flags: &FlagsState{Strict: true},
			Config: `
global {
	location = "` + escapedGlobal + `"
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"

	asset "secure" {
		pattern = "<img src=\"([^\"]+/secure/[^\"]+)"
		capture = 1
		find_all = true

		network {
			headers = {
				"custom_header" = "123"
			}
		}
	}
}`,
			Want: map[string]string{
				filepath.Join(global, "example", "a.jpg"): "securea",
				filepath.Join(global, "example", "b.jpg"): "secureb",
				filepath.Join(global, "example", "c.jpg"): "securec",
			},
			WantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(tc *testing.T) {
			// create a fresh filesystem
			utils.Fs, utils.AFS, utils.Wd = tu.SetupMemMapFs(root)

			mock := createMockGetCmd()
			g := New(mock)
			g.Flags = tt.Flags

			config, _, regexCache, diags := config.Parse([]byte(tt.Config), "test.hcl")
			if diags.HasErrors() {
				tc.Errorf("got errors: %+v", diags)
			}
			g.Config = config
			g.RegexCache = regexCache
			g.URLs = []string{ts.URL + "/gallery/123/test?id=543"}

			g.BuildSiteCache()

			cacheDiags := g.BuildAssetCache()
			if cacheDiags.HasErrors() {
				t.Fatalf("got errors: %+v", cacheDiags)
			}

			got := g.Download()

			if got.HasErrors() != tt.WantErr {
				t.Errorf("got %v, want errors %v", got, tt.WantErr)
			}

			for path, contents := range tt.Want {
				got, err := utils.AFS.ReadFile(path)
				if err != nil {
					t.Fatalf("got error: %v", err)
				}
				if string(got) != contents {
					t.Errorf("got %q, want %q", string(got), contents)
				}
			}
		})
	}
}
