package cmd

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/everdrone/grab/internal/utils"
	tu "github.com/everdrone/grab/testutils"
)

func TestGetCmd(t *testing.T) {
	initialWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(initialWd)
	}()

	root := tu.GetOSRoot()
	globalLocation := filepath.Join(root, "global")
	escapedGlobalLocation := tu.EscapeHCLString(globalLocation)

	// create test server
	e := tu.CreateMockServer()
	ts := httptest.NewUnstartedServer(e)

	ts.Listener.Close()
	ts.Listener = e.Listener
	ts.Start()

	defer ts.Close()

	tests := []struct {
		Name       string
		Args       []string
		Config     string
		ConfigPath string
		CheckFiles map[string]string
		WantErr    bool
	}{
		{
			Name: "invalid config",
			Args: []string{ts.URL + "/gallery/123/test"},
			Config: `
global {
	location = "` + escapedGlobalLocation + `"
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"
}
`,
			ConfigPath: filepath.Join(root, "grab.hcl"),
			WantErr:    true,
		},
		{
			Name: "invalid url",
			Args: []string{"1http://anything.com/fails"},
			Config: `
global {
	location = "` + escapedGlobalLocation + `"
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"

	asset "image" {
		pattern = "dummy"
		capture = 0
	}
}
`,
			ConfigPath: filepath.Join(root, "grab.hcl"),
			WantErr:    true,
		},
		{
			Name: "strict stops during build asset cache",
			Args: []string{ts.URL + "/givesNotFound", "-s"},
			Config: `
global {
	location = "` + escapedGlobalLocation + `"
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"

	asset "image" {
		pattern = "<img src=\"([^\"]+/img/[^\"]+)"
		capture = 1
		find_all = true
	}
}
`,
			ConfigPath: filepath.Join(root, "grab.hcl"),
			WantErr:    true,
		},
		{
			Name: "broken urls causes error",
			Args: []string{ts.URL + "/gallery/123/test", "-s"},
			Config: `
global {
	location = "` + escapedGlobalLocation + `"
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"

	asset "image" {
		pattern = "<a href=\"([^\"]*/broken/[^\"]+)"
		capture = 1
		find_all = true
	}
}
`,
			ConfigPath: filepath.Join(root, "grab.hcl"),
			WantErr:    true,
		},
		{
			Name: "can download",
			Args: []string{ts.URL + "/gallery/123/test", "-s"},
			Config: `
global {
	location = "` + escapedGlobalLocation + `"
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"

	asset "image" {
		pattern = "<img src=\"([^\"]+/img/[^\"]+)"
		capture = 1
		find_all = true
	}
}
`,
			CheckFiles: map[string]string{
				filepath.Join(globalLocation, "example", "a.jpg"): "imagea",
				filepath.Join(globalLocation, "example", "b.jpg"): "imageb",
				filepath.Join(globalLocation, "example", "c.jpg"): "imagec",
			},
			ConfigPath: filepath.Join(root, "grab.hcl"),
			WantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(tc *testing.T) {
			// reset fs
			utils.Fs, utils.AFS, utils.Wd = tu.SetupMemMapFs(root)

			// create config file and global dir
			utils.AFS.MkdirAll(globalLocation, os.ModePerm)
			utils.AFS.WriteFile(tt.ConfigPath, []byte(tt.Config), os.ModePerm)

			c, _, _, err := tu.ExecuteCommandErr(RootCmd, append([]string{"get"}, tt.Args...)...)

			if c.Name() != "get" {
				tc.Fatalf("got: %s, want: 'find", c.Name())
			}

			if tt.CheckFiles != nil {
				for f, v := range tt.CheckFiles {
					if got, _ := utils.AFS.ReadFile(f); string(got) != v {
						tc.Fatalf("got: %s, want: %s", string(got), v)
					}
				}
			}

			if (err != nil) != tt.WantErr {
				tc.Errorf("got: %v, want: %v", err, tt.WantErr)
			}
		})
	}
}
