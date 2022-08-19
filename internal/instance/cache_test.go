package instance

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/everdrone/grab/internal/config"
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
				"example\\.com": regexp.MustCompile(`example\.com`),
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
				"example\\.com": regexp.MustCompile(`example\.com`),
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
