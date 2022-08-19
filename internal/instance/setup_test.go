package instance

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"

	"github.com/everdrone/grab/internal/config"
	"github.com/everdrone/grab/internal/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

func createMockGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Just a mock",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().BoolP("force", "f", false, "overwrite existing files")
	cmd.Flags().StringP("config", "c", "", "the path of the config file to use")

	cmd.Flags().BoolP("strict", "s", false, "fail on errors")
	cmd.Flags().BoolP("dry-run", "n", false, "do not write on disk")

	cmd.Flags().BoolP("progress", "p", false, "show progress bars")
	cmd.Flags().BoolP("quiet", "q", false, "do not emit any output")
	cmd.Flags().CountP("verbose", "v", "verbosity level")

	return cmd
}

func TestNew(t *testing.T) {
	t.Run("new returns instance", func(tc *testing.T) {
		mock := createMockGetCmd()

		got := New(mock)
		if got == nil {
			tc.Errorf("got: %v, want: %v", got, nil)
		}

		if got.Command != mock {
			tc.Errorf("got: %v, want: %v", got.Command, mock)
		}
	})
}

func TestParseFlags(t *testing.T) {
	t.Run("parse quiet", func(tc *testing.T) {
		mock := createMockGetCmd()
		g := New(mock)

		mock.SetArgs([]string{"http://example.com", "-vvv", "-q"})
		mock.ExecuteC()

		g.ParseFlags()

		if g.Flags.Quiet != true {
			tc.Errorf("got: %v, want: %v", g.Flags.Quiet, true)
		}
		if g.Flags.Verbosity != 0 {
			tc.Errorf("got: %v, want: %v", g.Flags.Verbosity, 0)
		}
	})

	t.Run("parse verbosity", func(tc *testing.T) {
		mock := createMockGetCmd()
		g := New(mock)

		mock.SetArgs([]string{"http://example.com", "-vv"})
		mock.ExecuteC()

		g.ParseFlags()

		if g.Flags.Verbosity != 3 {
			tc.Errorf("got: %v, want: %v", g.Flags.Verbosity, 3)
		}
	})
}

func TestParseConfig(t *testing.T) {
	initialWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(initialWd)
	}()

	root := utils.GetOSRoot()

	homedir, _ := homedir.Dir()

	tests := []struct {
		Name           string
		Flags          *FlagsState
		Wd             string
		ConfigFilePath string
		Config         string
		WantConfig     *config.Config
		WantRegexCache config.RegexCacheMap
		WantErr        bool
	}{
		{
			Name:           "config ok",
			Flags:          &FlagsState{},
			Wd:             filepath.Join(root, "test"),
			ConfigFilePath: filepath.Join(root, "test", "grab.hcl"),
			Config: `
global {
	location = "` + utils.EscapeHCLString(filepath.Join(root, "home", "user", "downloads")) + `"
}

site "example" {
	test = "testPattern"

	asset "image" {
		pattern = "assetPattern"
		capture = 0
	}
}`,
			WantConfig: &config.Config{
				Global: config.GlobalConfig{
					Location: filepath.Join(root, "home", "user", "downloads"),
				},
				Sites: []config.SiteConfig{
					{
						Name: "example",
						Test: "testPattern",
						Assets: []config.AssetConfig{
							{
								Name:    "image",
								Pattern: "assetPattern",
								Capture: "0",
							},
						},
					},
				},
			},
			WantRegexCache: config.RegexCacheMap{
				"assetPattern": regexp.MustCompile("assetPattern"),
				"testPattern":  regexp.MustCompile("testPattern"),
			},
			WantErr: false,
		},
		{
			Name:           "expands home directory",
			Flags:          &FlagsState{},
			Wd:             filepath.Join(root, "test"),
			ConfigFilePath: filepath.Join(root, "test", "grab.hcl"),
			Config: `
global {
	location = "` + utils.EscapeHCLString(filepath.Join("~", "Downloads", "grab")) + `"
}

site "example" {
	test = "testPattern"

	asset "image" {
		pattern = "assetPattern"
		capture = 0
	}
}`,
			WantConfig: &config.Config{
				Global: config.GlobalConfig{
					Location: filepath.Join(homedir, "Downloads", "grab"),
				},
				Sites: []config.SiteConfig{
					{
						Name: "example",
						Test: "testPattern",
						Assets: []config.AssetConfig{
							{
								Name:    "image",
								Pattern: "assetPattern",
								Capture: "0",
							},
						},
					},
				},
			},
			WantRegexCache: config.RegexCacheMap{
				"assetPattern": regexp.MustCompile("assetPattern"),
				"testPattern":  regexp.MustCompile("testPattern"),
			},
			WantErr: false,
		},
		{
			Name:           "expands relative path",
			Flags:          &FlagsState{},
			Wd:             filepath.Join(root, "test"),
			ConfigFilePath: filepath.Join(root, "test", "grab.hcl"),
			Config: `
global {
	location = "` + utils.EscapeHCLString(filepath.Join("..", "expandMe")) + `"
}

site "example" {
	test = "testPattern"

	asset "image" {
		pattern = "assetPattern"
		capture = 0
	}
}`,
			WantConfig: &config.Config{
				Global: config.GlobalConfig{
					Location: filepath.Join(root, "expandMe"),
				},
				Sites: []config.SiteConfig{
					{
						Name: "example",
						Test: "testPattern",
						Assets: []config.AssetConfig{
							{
								Name:    "image",
								Pattern: "assetPattern",
								Capture: "0",
							},
						},
					},
				},
			},
			WantRegexCache: config.RegexCacheMap{
				"assetPattern": regexp.MustCompile("assetPattern"),
				"testPattern":  regexp.MustCompile("testPattern"),
			},
			WantErr: false,
		},
		{
			Name:           "config not found",
			Flags:          &FlagsState{},
			Wd:             filepath.Join(root, "test"),
			ConfigFilePath: filepath.Join(root, "test", "deeper", "grab.hcl"),
			Config:         "",
			WantConfig:     nil,
			WantRegexCache: config.RegexCacheMap(nil),
			WantErr:        true,
		},
		{
			Name:           "invalid config",
			Flags:          &FlagsState{},
			Wd:             filepath.Join(root, "test"),
			ConfigFilePath: filepath.Join(root, "test", "grab.hcl"),
			Config: `
site "example" {
	test = "testPattern"

	asset "image" {
		pattern = "assetPattern"
		capture = 0
	}
}`,
			WantConfig:     nil,
			WantRegexCache: config.RegexCacheMap(nil),
			WantErr:        true,
		},
		{
			Name:           "cannot expand home directory",
			Flags:          &FlagsState{},
			Wd:             filepath.Join(root, "test"),
			ConfigFilePath: filepath.Join(root, "test", "grab.hcl"),
			Config: `
global {
	location = "` + filepath.Join("~user", "Downloads", "grab") + `"
}

site "example" {
	test = "testPattern"

	asset "image" {
		pattern = "assetPattern"
		capture = 0
	}
}`,
			WantConfig:     nil,
			WantRegexCache: config.RegexCacheMap(nil),
			WantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(tc *testing.T) {
			// create file system
			utils.SetupMemMapFs(root)
			utils.AFS.WriteFile(tt.ConfigFilePath, []byte(tt.Config), os.ModePerm)

			// set working directory
			func() {
				utils.Wd = tt.Wd
			}()

			mock := createMockGetCmd()
			g := New(mock)

			// set flags
			g.Flags = tt.Flags

			diags := g.ParseConfig()

			if diags.HasErrors() != tt.WantErr {
				tc.Errorf("got: %v, want: %v", diags, tt.WantErr)
			}

			// do not check other outputs if we got errors
			if !diags.HasErrors() {
				if !reflect.DeepEqual(g.Config, tt.WantConfig) {
					tc.Errorf("got: %+v, want: %+v", g.Config, tt.WantConfig)
				}

				if !reflect.DeepEqual(g.RegexCache, tt.WantRegexCache) {
					tc.Errorf("got: %+v, want: %+v", g.RegexCache, tt.WantRegexCache)
				}
			}
		})
	}
}
