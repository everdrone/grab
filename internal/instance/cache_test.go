package instance

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"

	"github.com/everdrone/grab/internal/config"
	"github.com/everdrone/grab/internal/utils"
	"github.com/labstack/echo/v4"
)

func TestBuildSiteCache(t *testing.T) {
	tests := []struct {
		Name       string
		RegexCache config.RegexCacheMap
		URLs       []string
		Config     *config.Config
		Want       *config.Config
	}{
		{
			Name:       "no urls",
			Config:     &config.Config{},
			RegexCache: config.RegexCacheMap{},
			URLs:       []string{},
			Want:       &config.Config{},
		},
		{
			Name: "one url",
			RegexCache: config.RegexCacheMap{
				`example\.com`: regexp.MustCompile(`example\.com`),
			},
			URLs: []string{"https://example.com/gallery/test"},
			Config: &config.Config{
				Sites: []config.SiteConfig{
					{Name: "example",
						Test: "example\\.com",
					},
				},
			},
			Want: &config.Config{
				Sites: []config.SiteConfig{
					{Name: "example",
						Test: "example\\.com",
						URLs: []string{"https://example.com/gallery/test"},
					},
				},
			},
		},
		{
			Name: "multiple urls",
			RegexCache: config.RegexCacheMap{
				`example\.com`: regexp.MustCompile(`example\.com`),
			},
			URLs: []string{"https://example.com/gallery/test", "https://example.com/other"},
			Config: &config.Config{
				Sites: []config.SiteConfig{
					{Name: "example",
						Test: "example\\.com",
					},
				},
			},
			Want: &config.Config{
				Sites: []config.SiteConfig{
					{Name: "example",
						Test: "example\\.com",
						URLs: []string{"https://example.com/gallery/test", "https://example.com/other"},
					},
				},
			},
		},
		{
			Name: "no matches",
			RegexCache: config.RegexCacheMap{
				"example\\.com": regexp.MustCompile(`example\.com`),
			},
			URLs: []string{"https://not-matching.com/gallery/1", "https://not-matching.com/gallery/1"},
			Config: &config.Config{
				Sites: []config.SiteConfig{
					{Name: "example",
						Test: "example\\.com",
					},
				},
			},
			Want: &config.Config{
				Sites: []config.SiteConfig{
					{Name: "example",
						Test: "example\\.com",
						URLs: []string(nil),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(tc *testing.T) {
			g := New(nil)

			g.Config = tt.Config
			g.RegexCache = tt.RegexCache
			g.URLs = tt.URLs

			g.BuildSiteCache()

			if !reflect.DeepEqual(g.Config, tt.Want) {
				tc.Errorf("got: %+v, want: %+v", g.Config, tt.Want)
			}
		})
	}
}

func TestBuildAssetCache(t *testing.T) {
	// FIXME: let's create a function and set up an in-memory file server for testing purposes
	testPath := `/gallery/123/test?id=123`
	testPage := `example page
https://example.com/img/test1.jpg
https://example.com/img/test2.jpg
https://example.com/img/test3.jpg

https://example.com/video/123abc/small.mp4
https://example.com/video/def456/small.mp4
https://example.com/video/ghi789/small.mp4

<img src="/relative/image.jpg" />

name: foo
description: bar
`

	root := utils.GetOSRoot()
	globalLocation := filepath.Join(root, "global")

	// create test server
	e := echo.New()
	e.GET("/gallery/:id/:name", func(c echo.Context) error {
		return c.HTML(http.StatusOK, testPage)
	})
	ts := httptest.NewServer(e)
	defer ts.Close()

	tests := []struct {
		Name    string
		Flags   *FlagsState
		URLs    []string
		Config  string
		Want    *config.Config
		WantErr bool
	}{
		{
			Name:  "not found",
			Flags: &FlagsState{},
			URLs:  []string{ts.URL + "/notFound"},
			Config: `
global {
	location = "` + utils.EscapeHCLString(globalLocation) + `"
}

site "example" {
	test = "http://(127\\.0\\.0\\.1|localhost):"
	asset "image" {
		pattern = "https:\\/\\/example\\.com\\/img\\/\\w+\\.\\w+"
		capture = 0
		find_all = true
	}
}`,
			Want: &config.Config{
				Sites: []config.SiteConfig{
					{
						Assets: []config.AssetConfig{
							{
								Downloads: nil,
							},
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name:  "not found strict",
			Flags: &FlagsState{Strict: true},
			URLs:  []string{ts.URL + "/notFound"},
			Config: `
global {
	location = "` + utils.EscapeHCLString(globalLocation) + `"
}

site "example" {
	test = "http://(127\\.0\\.0\\.1|localhost):"
	asset "image" {
		pattern = "https:\\/\\/example\\.com\\/img\\/\\w+\\.\\w+"
		capture = 0
		find_all = true
	}
}`,
			Want: &config.Config{
				Sites: []config.SiteConfig{
					{
						Assets: []config.AssetConfig{
							{
								Downloads: nil,
							},
						},
					},
				},
			},
			WantErr: true,
		},
		{
			Name:  "no urls",
			Flags: &FlagsState{},
			Config: `
global {
	location = "` + utils.EscapeHCLString(globalLocation) + `"
}

site "example" {
	test = "http://(127\\.0\\.0\\.1|localhost):"
	asset "image" {
		pattern = "https:\\/\\/example\\.com\\/img\\/\\w+\\.\\w+"
		capture = 0
		find_all = true
	}
}`,
			Want: &config.Config{
				Sites: []config.SiteConfig{
					{
						Assets: []config.AssetConfig{
							{
								Downloads: nil,
							},
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name:  "simple assets",
			Flags: &FlagsState{},
			URLs:  []string{ts.URL + testPath},
			Config: `
global {
	location = "` + utils.EscapeHCLString(globalLocation) + `"
}

site "example" {
	test = "http://(127\\.0\\.0\\.1|localhost):"
	asset "image" {
		pattern = "https:\\/\\/example\\.com\\/img\\/\\w+\\.\\w+"
		capture = 0
		find_all = true
	}
}`,
			Want: &config.Config{
				Sites: []config.SiteConfig{
					{
						Assets: []config.AssetConfig{
							{
								Downloads: map[string]string{
									"https://example.com/img/test1.jpg": filepath.Join(globalLocation, "example", "test1.jpg"),
									"https://example.com/img/test2.jpg": filepath.Join(globalLocation, "example", "test2.jpg"),
									"https://example.com/img/test3.jpg": filepath.Join(globalLocation, "example", "test3.jpg"),
								},
							},
						},
						InfoMap: map[string]map[string]string{
							filepath.Join(globalLocation, "example"): {
								"url": ts.URL + testPath,
							},
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name:  "subdirectory",
			Flags: &FlagsState{},
			URLs:  []string{ts.URL + testPath},
			Config: `
global {
	location = "` + utils.EscapeHCLString(globalLocation) + `"
}

site "example" {
	test = "http://(127\\.0\\.0\\.1|localhost):"
	asset "image" {
		pattern = "https:\\/\\/example\\.com\\/img\\/\\w+\\.\\w+"
		capture = 0
		find_all = true
	}

	subdirectory {
		pattern = "\\/gallery\\/(\\d+)"
		capture = 1
		from = url
	}
}`,
			Want: &config.Config{
				Sites: []config.SiteConfig{
					{
						Assets: []config.AssetConfig{
							{
								Downloads: map[string]string{
									"https://example.com/img/test1.jpg": filepath.Join(globalLocation, "example", "123", "test1.jpg"),
									"https://example.com/img/test2.jpg": filepath.Join(globalLocation, "example", "123", "test2.jpg"),
									"https://example.com/img/test3.jpg": filepath.Join(globalLocation, "example", "123", "test3.jpg"),
								},
							},
						},
						InfoMap: map[string]map[string]string{
							filepath.Join(globalLocation, "example", "123"): {
								"url": ts.URL + testPath,
							},
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name:  "absolute subdirectory",
			Flags: &FlagsState{},
			URLs:  []string{ts.URL + testPath},
			Config: `
global {
	location = "` + utils.EscapeHCLString(globalLocation) + `"
}

site "example" {
	test = "http://(127\\.0\\.0\\.1|localhost):"
	asset "image" {
		pattern = "https:\\/\\/example\\.com\\/img\\/\\w+\\.\\w+"
		capture = 0
		find_all = true
	}

	subdirectory {
		pattern = "\\/gallery\\/(\\d+)"
		capture = 1
		from = url
	}
}`,
			Want: &config.Config{
				Sites: []config.SiteConfig{
					{
						Assets: []config.AssetConfig{
							{
								Downloads: map[string]string{
									"https://example.com/img/test1.jpg": filepath.Join(globalLocation, "example", "123", "test1.jpg"),
									"https://example.com/img/test2.jpg": filepath.Join(globalLocation, "example", "123", "test2.jpg"),
									"https://example.com/img/test3.jpg": filepath.Join(globalLocation, "example", "123", "test3.jpg"),
								},
							},
						},
						InfoMap: map[string]map[string]string{
							filepath.Join(globalLocation, "example", "123"): {
								"url": ts.URL + testPath,
							},
						},
					},
				},
			},
			WantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(tc *testing.T) {
			g := New(nil)
			g.Flags = tt.Flags

			config, _, regexCache, diags := config.Parse([]byte(tt.Config), "test.hcl")
			if diags.HasErrors() {
				tc.Errorf("got errors: %+v", diags)
			}
			g.Config = config
			g.RegexCache = regexCache

			g.URLs = tt.URLs
			g.BuildSiteCache()

			gotDiags := g.BuildAssetCache()

			if gotDiags.HasErrors() != tt.WantErr {
				tc.Errorf("got: %+v, want errors: %+v", gotDiags.HasErrors(), tt.WantErr)
			}

			for i, site := range g.Config.Sites {
				for j, asset := range site.Assets {
					if !reflect.DeepEqual(asset.Downloads, tt.Want.Sites[i].Assets[j].Downloads) {
						tc.Errorf("got: %+v, want: %+v", asset.Downloads, tt.Want.Sites[i].Assets[j].Downloads)
					}
				}

				gotInfoMap := site.InfoMap
				wantInfoMap := tt.Want.Sites[i].InfoMap

				compareInfoMaps(tc, gotInfoMap, wantInfoMap)
			}
		})
	}
}

func compareInfoMaps(t *testing.T, got, want config.InfoCacheMap) {
	// check that the keys are the same
	if !reflect.DeepEqual(getMapKeys(got), getMapKeys(want)) {
		t.Errorf("got: %+v, want: %+v", getMapKeys(got), getMapKeys(want))
	}

	// check that the values are the same but ignore the timestamp
	for k, v := range got {
		for k2, v2 := range v {
			if k2 == "timestamp" {
				continue
			}
			if want[k][k2] != v2 {
				t.Errorf("got: %+v, want: %+v", got, want)
			}
		}
	}
}

func getMapKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
