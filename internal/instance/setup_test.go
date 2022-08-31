package instance

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"

	"github.com/everdrone/grab/internal/config"
	"github.com/everdrone/grab/internal/utils"
	tu "github.com/everdrone/grab/testutils"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
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
	tests := []struct {
		Name            string
		Args            []string
		Want            *FlagsState
		WantGlobalLevel zerolog.Level
	}{
		{
			"quiet",
			[]string{"-q"},
			&FlagsState{
				Quiet:     true,
				Verbosity: 0,
			},
			zerolog.ErrorLevel,
		},
		{
			"quiet and verbosity",
			[]string{"-q", "-vvv"},
			&FlagsState{
				Quiet:     true,
				Verbosity: 0,
			},
			zerolog.ErrorLevel,
		},
		{
			"force",
			[]string{"-f"},
			&FlagsState{
				Force:     true,
				Verbosity: 1,
			},
			zerolog.WarnLevel,
		},
		{
			"zerolog info level",
			[]string{"-v"},
			&FlagsState{
				Verbosity: 2,
			},
			zerolog.InfoLevel,
		},
		{
			"zerolog debug level",
			[]string{"-vv"},
			&FlagsState{
				Verbosity: 3,
			},
			zerolog.DebugLevel,
		},
		{
			"zerolog trace level",
			[]string{"-vvv"},
			&FlagsState{
				Verbosity: 4,
			},
			zerolog.TraceLevel,
		},
		{
			"too many verbosity levels",
			[]string{"-vvvvv"},
			&FlagsState{
				Verbosity: 6,
			},
			zerolog.TraceLevel,
		},
		{
			"config path",
			[]string{"-c", "grab.hcl"},
			&FlagsState{
				ConfigPath: "grab.hcl",
				Verbosity:  1,
			},
			zerolog.WarnLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(tc *testing.T) {
			mock := createMockGetCmd()
			mock.SetArgs(append([]string{"list.ini"}, tt.Args...))
			mock.ExecuteC()
			g := New(mock)

			g.ParseFlags()

			if !reflect.DeepEqual(g.Flags, tt.Want) {
				tc.Errorf("got: %+v, want: %+v", g.Flags, tt.Want)
			}

		})
	}
}

func TestParseConfig(t *testing.T) {
	initialWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(initialWd)
	}()

	root := tu.GetOSRoot()
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
	location = "` + tu.EscapeHCLString(filepath.Join(root, "home", "user", "downloads")) + `"
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
	location = "` + tu.EscapeHCLString(filepath.Join("~", "Downloads", "grab")) + `"
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
	location = "` + tu.EscapeHCLString(filepath.Join("..", "expandMe")) + `"
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
	location = "` + tu.EscapeHCLString(filepath.Join("~user", "Downloads", "grab")) + `"
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
		{
			Name: "cannot read config file",
			Flags: &FlagsState{
				ConfigPath: filepath.Join(root, "test", "file_not_readable.txt"),
			},
			Wd:             filepath.Join(root, "test"),
			ConfigFilePath: filepath.Join(root, "test", "file_not_readable.txt"),
			Config: `
global {
	location = "` + tu.EscapeHCLString(filepath.Join("~user", "Downloads", "grab")) + `"
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
			utils.Fs, utils.AFS, utils.Wd = tu.SetupMemMapFs(root)
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
				tc.Errorf("got: %v, want errors: %v", diags, tt.WantErr)
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

func TestParseURLs(t *testing.T) {
	initialWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(initialWd)
	}()
	root := tu.GetOSRoot()

	// create file system
	utils.Fs, utils.AFS, utils.Wd = tu.SetupMemMapFs(root)
	utils.AFS.MkdirAll(filepath.Join(root, "test"), os.ModePerm)
	utils.AFS.WriteFile(filepath.Join(root, "test", "list.ini"), []byte(`# https://ignore.me.com/image1
https://example.com/gallery/test
   http://example.com/g2
; ignore this as well
`), os.ModePerm)
	utils.AFS.WriteFile(filepath.Join(root, "invalid.ini"), []byte(`; comment
1https://example.com/gallery/test
# comment 2
`), os.ModePerm)

	tests := []struct {
		Name     string
		Wd       string
		Args     []string
		WantURLs []string
		WantErr  bool
	}{
		{
			// this case should not happen since we require at least one argument
			Name:     "no args",
			Wd:       root,
			Args:     []string(nil),
			WantURLs: []string(nil),
			WantErr:  false,
		},
		{
			Name:     "one url",
			Wd:       root,
			Args:     []string{"https://example.com/gallery/test"},
			WantURLs: []string{"https://example.com/gallery/test"},
			WantErr:  false,
		},
		{
			Name:     "more urls",
			Wd:       root,
			Args:     []string{"https://example.com/gallery/test", "https://example.com/gallery/test2"},
			WantURLs: []string{"https://example.com/gallery/test", "https://example.com/gallery/test2"},
			WantErr:  false,
		},
		{
			Name:     "duplicate urls",
			Wd:       root,
			Args:     []string{"https://example.com/gallery/test", "https://example.com/gallery/test"},
			WantURLs: []string{"https://example.com/gallery/test"},
			WantErr:  false,
		},
		{
			Name:     "one file",
			Wd:       root,
			Args:     []string{filepath.Join(root, "test", "list.ini")},
			WantURLs: []string{"https://example.com/gallery/test", "http://example.com/g2"},
			WantErr:  false,
		},
		{
			Name:     "duplicate files",
			Wd:       root,
			Args:     []string{filepath.Join(root, "test", "list.ini"), filepath.Join(root, "test", "list.ini")},
			WantURLs: []string{"https://example.com/gallery/test", "http://example.com/g2"},
			WantErr:  false,
		},
		{
			Name:     "urls and files",
			Wd:       root,
			Args:     []string{"http://example.com/1", filepath.Join(root, "test", "list.ini")},
			WantURLs: []string{"http://example.com/1", "https://example.com/gallery/test", "http://example.com/g2"},
			WantErr:  false,
		},
		{
			Name:     "url is also in file",
			Wd:       root,
			Args:     []string{"https://example.com/gallery/test", filepath.Join(root, "test", "list.ini")},
			WantURLs: []string{"https://example.com/gallery/test", "http://example.com/g2"},
			WantErr:  false,
		},
		{
			Name:     "invalid url in file",
			Wd:       root,
			Args:     []string{filepath.Join(root, "invalid.ini")},
			WantURLs: []string(nil),
			WantErr:  true,
		},
		{
			Name:     "file not found",
			Wd:       root,
			Args:     []string{filepath.Join(root, "notFound.ini")},
			WantURLs: []string(nil),
			WantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(tc *testing.T) {
			func() {
				utils.Wd = tt.Wd
			}()

			mock := createMockGetCmd()
			g := New(mock)

			diags := g.ParseURLs(tt.Args)
			if diags.HasErrors() != tt.WantErr {
				tc.Errorf("got: %v, want errors: %v", diags, tt.WantErr)
			}
			if !diags.HasErrors() {
				if !reflect.DeepEqual(g.URLs, tt.WantURLs) {
					tc.Errorf("got: %+v, want: %+v", g.URLs, tt.WantURLs)
				}
			}
		})
	}
}
