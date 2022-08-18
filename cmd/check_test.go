package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/everdrone/grab/internal/utils"
)

const fileOk = `
global {
	location = "/home/user/Downloads/grab"
}

site "example" {
	test = ":\\/\\/example\\.com"

	asset "image" {
		pattern  = "<img\\ssrc=\"([^\"]+)"
		capture  = 1
		find_all = true
	}

	asset "video" {
		pattern  = "<video\\ssrc=\"(?P<video_url>[^\"]+)"
		capture  = "video_url"
		find_all = true
	}
}
`

const fileInvalid = `
global {
	location = "/home/user/Downloads/grab"
}
`

func TestCheckCmd(t *testing.T) {
	initialWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(initialWd)
	}()

	root := utils.GetOSRoot()
	utils.SetupMemMapFs(root)

	utils.Fs.MkdirAll("/other/directory", os.ModePerm)
	utils.Fs.MkdirAll("/tmp/test/config/nested", os.ModePerm)
	utils.AFS.WriteFile("/tmp/test/config/grab.hcl", []byte(fileOk), os.ModePerm)
	utils.AFS.WriteFile("/tmp/test/config/invalid.hcl", []byte(fileInvalid), os.ModePerm)

	tests := []struct {
		Name         string
		Wd           string
		Args         []string
		WantContains string
		WantErr      string
		HasErrors    bool
	}{
		{
			Name:         "ok nested",
			Wd:           "/tmp/test/config/nested",
			Args:         []string{},
			WantContains: "ok\n",
			WantErr:      "",
			HasErrors:    false,
		},
		{
			Name:         "ok same directory",
			Wd:           "/tmp/test/config",
			Args:         []string{},
			WantContains: "ok\n",
			WantErr:      "",
			HasErrors:    false,
		},
		{
			Name:         "ok quiet",
			Wd:           "/tmp/test/config/nested",
			Args:         []string{"-q"},
			WantContains: "",
			WantErr:      "",
			HasErrors:    false,
		},
		{
			Name:         "not found",
			Wd:           "/other/directory",
			Args:         []string{},
			WantContains: "",
			WantErr: `╷ Error: could not resolve config file
╵   could not resolve grab.hcl
`,
			HasErrors: true,
		},
		{
			Name:         "invalid args",
			Wd:           "/tmp/test/config/nested",
			Args:         []string{"--unknown"},
			WantContains: "Usage:",
			WantErr:      "Error: unknown flag: --unknown\n",
			HasErrors:    true,
		},
		{
			Name:         "missing global",
			Wd:           "/tmp/test/config/nested",
			Args:         []string{"-c", "/tmp/test/config/invalid.hcl"},
			WantContains: "",
			WantErr: `╷ Error: Insufficient site blocks
│   At least 1 "site" blocks are required.
╵   /tmp/test/config/invalid.hcl:1,1-1
`,
			HasErrors: true,
		},
		{
			Name:         "file does not exist",
			Wd:           "/tmp/test/config/nested",
			Args:         []string{"-c", "notFound.hcl"},
			WantContains: "",
			WantErr: `╷ Error: could not read config file
╵   open notFound.hcl:`,
			HasErrors: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			args := []string{"config", "check"}

			func() {
				utils.Wd = tt.Wd
			}()

			c, got, gotErr, err := utils.ExecuteCommandErr(RootCmd, append(args, tt.Args...)...)
			if err != nil && !tt.HasErrors {
				t.Errorf("unexpected error: %v", err)
			}

			if c.Name() != "check" {
				t.Errorf("got: '%s', want: 'check'", c.Name())
			}
			if !strings.Contains(got, tt.WantContains) {
				t.Errorf("got: %s, does not contain: %s", got, tt.WantContains)
			}
			if !strings.HasPrefix(gotErr, tt.WantErr) {
				t.Errorf("got: %s, want: %s", gotErr, tt.WantErr)
			}
		})
	}
}
