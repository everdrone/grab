package instance

import (
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"

	"github.com/everdrone/grab/internal/config"
	"github.com/everdrone/grab/internal/testutils"
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

// FIXME: from this line down, the code is a mess.
// it does test the functionality of cache.go but it's very very messy.
// it should be refactored.
func TestBuildAssetCache(t *testing.T) {
	testPath := `/gallery/123/test?id=543`
	root := testutils.GetOSRoot()
	globalLocation := filepath.Join(root, "global")

	// create test server
	e := testutils.CreateMockServer()

	// hacky way of getting the same port as echo's listener
	// see: https://stackoverflow.com/a/42218765
	ts := httptest.NewUnstartedServer(e)

	ts.Listener.Close()
	ts.Listener = e.Listener
	ts.Start()

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
	location = "` + testutils.EscapeHCLString(globalLocation) + `"
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
	location = "` + testutils.EscapeHCLString(globalLocation) + `"
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
	location = "` + testutils.EscapeHCLString(globalLocation) + `"
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
			Name:  "find one asset",
			Flags: &FlagsState{},
			URLs:  []string{ts.URL + testPath},
			Config: `
global {
	location = "` + testutils.EscapeHCLString(globalLocation) + `"
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"
	asset "image" {
		pattern = "<img src=\"([^\"]+/img/[^\"]+)"
		capture = 1
		find_all = false
	}
}`,
			Want: &config.Config{
				Sites: []config.SiteConfig{
					{
						Assets: []config.AssetConfig{
							{
								Downloads: map[string]string{
									ts.URL + "/img/a.jpg": filepath.Join(globalLocation, "example", "a.jpg"),
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
			Name:  "simple assets",
			Flags: &FlagsState{},
			URLs:  []string{ts.URL + testPath},
			Config: `
global {
	location = "` + testutils.EscapeHCLString(globalLocation) + `"
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"
	asset "image" {
		pattern = "<img src=\"([^\"]+/img/[^\"]+)"
		capture = 1
		find_all = true
	}
}`,
			Want: &config.Config{
				Sites: []config.SiteConfig{
					{
						Assets: []config.AssetConfig{
							{
								Downloads: map[string]string{
									ts.URL + "/img/a.jpg": filepath.Join(globalLocation, "example", "a.jpg"),
									ts.URL + "/img/b.jpg": filepath.Join(globalLocation, "example", "b.jpg"),
									ts.URL + "/img/c.jpg": filepath.Join(globalLocation, "example", "c.jpg"),
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
			Name:  "relative assets",
			Flags: &FlagsState{},
			URLs:  []string{ts.URL + testPath},
			Config: `
global {
	location = "` + testutils.EscapeHCLString(globalLocation) + `"
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"
	asset "image" {
		pattern = "<img src=\"(/img/[^\"]+)"
		capture = 1
		find_all = true
	}
}`,
			Want: &config.Config{
				Sites: []config.SiteConfig{
					{
						Assets: []config.AssetConfig{
							{
								Downloads: map[string]string{
									ts.URL + "/img/a.jpg": filepath.Join(globalLocation, "example", "a.jpg"),
									ts.URL + "/img/b.jpg": filepath.Join(globalLocation, "example", "b.jpg"),
									ts.URL + "/img/c.jpg": filepath.Join(globalLocation, "example", "c.jpg"),
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
			Name:  "subdirectory from url",
			Flags: &FlagsState{},
			URLs:  []string{ts.URL + testPath},
			Config: `
global {
	location = "` + testutils.EscapeHCLString(globalLocation) + `"
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"
	asset "image" {
		pattern = "<img src=\"([^\"]+/img/[^\"]+)"
		capture = 1
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
									ts.URL + "/img/a.jpg": filepath.Join(globalLocation, "example", "123", "a.jpg"),
									ts.URL + "/img/b.jpg": filepath.Join(globalLocation, "example", "123", "b.jpg"),
									ts.URL + "/img/c.jpg": filepath.Join(globalLocation, "example", "123", "c.jpg"),
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
			Name:  "subdirectory from body",
			Flags: &FlagsState{},
			URLs:  []string{ts.URL + testPath},
			Config: `
global {
	location = "` + testutils.EscapeHCLString(globalLocation) + `"
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"
	asset "image" {
		pattern = "<img src=\"([^\"]+/img/[^\"]+)"
		capture = 1
		find_all = true
	}

	subdirectory {
		pattern = "Author: @(?P<username>[^<]+)"
		capture = "username"
		from = body
	}
}`,
			Want: &config.Config{
				Sites: []config.SiteConfig{
					{
						Assets: []config.AssetConfig{
							{
								Downloads: map[string]string{
									ts.URL + "/img/a.jpg": filepath.Join(globalLocation, "example", "everdrone", "a.jpg"),
									ts.URL + "/img/b.jpg": filepath.Join(globalLocation, "example", "everdrone", "b.jpg"),
									ts.URL + "/img/c.jpg": filepath.Join(globalLocation, "example", "everdrone", "c.jpg"),
								},
							},
						},
						InfoMap: map[string]map[string]string{
							filepath.Join(globalLocation, "example", "everdrone"): {
								"url": ts.URL + testPath,
							},
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name:  "transform url",
			Flags: &FlagsState{},
			URLs:  []string{ts.URL + testPath},
			Config: `
global {
	location = "` + testutils.EscapeHCLString(globalLocation) + `"
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"
	asset "video" {
		pattern = "<video src=\"([^\"]+)"
		capture = 1
		find_all = true

		transform url {
			pattern = "(.+)small(.*)"
			replace = "$${1}large$2"
		}
	}
}`,
			Want: &config.Config{
				Sites: []config.SiteConfig{
					{
						Assets: []config.AssetConfig{
							{
								Downloads: map[string]string{
									ts.URL + "/video/a/large.mp4": filepath.Join(globalLocation, "example", "large.mp4"),
									ts.URL + "/video/b/large.mp4": filepath.Join(globalLocation, "example", "large.mp4"),
									ts.URL + "/video/c/large.mp4": filepath.Join(globalLocation, "example", "large.mp4"),
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
			Name:  "transform filename",
			Flags: &FlagsState{},
			URLs:  []string{ts.URL + testPath},
			Config: `
global {
	location = "` + testutils.EscapeHCLString(globalLocation) + `"
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"
	asset "video" {
		pattern = "<video src=\"([^\"]+)"
		capture = 1
		find_all = true

		transform filename {
			pattern = ".+\\/video\\/(?P<id>\\w+)\\/(\\w+)\\.(?P<extension>\\w+)"
			replace = "$${id}.$${extension}"
		}
	}
}`,
			Want: &config.Config{
				Sites: []config.SiteConfig{
					{
						Assets: []config.AssetConfig{
							{
								Downloads: map[string]string{
									ts.URL + "/video/a/small.mp4": filepath.Join(globalLocation, "example", "a.mp4"),
									ts.URL + "/video/b/small.mp4": filepath.Join(globalLocation, "example", "b.mp4"),
									ts.URL + "/video/c/small.mp4": filepath.Join(globalLocation, "example", "c.mp4"),
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
			Name:  "transform filename absolute",
			Flags: &FlagsState{},
			URLs:  []string{ts.URL + testPath},
			Config: `
global {
	location = "` + testutils.EscapeHCLString(globalLocation) + `"
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"
	asset "video" {
		pattern = "<video src=\"([^\"]+)"
		capture = 1
		find_all = true

		transform filename {
			pattern = ".+\\/video\\/(?P<id>\\w+)\\/(\\w+)\\.(?P<extension>\\w+)"
			replace = "` + testutils.EscapeHCLString(root) + `$${id}.$${extension}"
		}
	}
}`,
			Want: &config.Config{
				Sites: []config.SiteConfig{
					{
						Assets: []config.AssetConfig{
							{
								Downloads: map[string]string{
									ts.URL + "/video/a/small.mp4": filepath.Join(root, "a.mp4"),
									ts.URL + "/video/b/small.mp4": filepath.Join(root, "b.mp4"),
									ts.URL + "/video/c/small.mp4": filepath.Join(root, "c.mp4"),
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
			Name:  "info",
			Flags: &FlagsState{},
			URLs:  []string{ts.URL + testPath},
			Config: `
global {
	location = "` + testutils.EscapeHCLString(globalLocation) + `"
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"
	info "author" {
		pattern = "Author: @(?P<username>[^<]+)"
		capture = "username"
	}

	info "title" {
		pattern = "<title>([^<]+)"
		capture = 1
	}
}`,
			Want: &config.Config{
				Sites: []config.SiteConfig{
					{
						Assets: []config.AssetConfig{
							{
								Downloads: map[string]string(nil),
							},
						},
						InfoMap: map[string]map[string]string{
							filepath.Join(globalLocation, "example"): {
								"url":    ts.URL + testPath,
								"author": "everdrone",
								"title":  "Grab Test Server",
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
