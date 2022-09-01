package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	tu "github.com/everdrone/grab/testutils"
	"github.com/hashicorp/hcl/v2"
	"golang.org/x/exp/slices"
)

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		URL  string
		Want bool
		Name string
	}{
		{URL: "", Want: false, Name: "empty"},
		{URL: "/foo/bar", Want: false, Name: "unix absolute path"},
		{URL: "://foo/bar", Want: false, Name: "no scheme"},
		{URL: "https://foo/bar", Want: true, Name: "no dot com ssl"},
		{URL: "tcp://foo/bar", Want: true, Name: "tcp no dot com"},
		{URL: "https://foo.com/bar", Want: true, Name: "valid ssl"},
		{URL: "1http://anything.com/fails", Want: false, Name: "invalid scheme"},
		{URL: "tcp://foo.com/bar", Want: true, Name: "valid tcp"},
		{URL: "c:\\windows\\bad", Want: false, Name: "windows absolute"},
		{URL: "\\unix\\good", Want: false, Name: "windows absolute without drive"},
		{URL: "unix\\good", Want: false, Name: "windows relative"},
		{URL: "foo/bar", Want: false, Name: "unix relative"},
		{URL: "../foo/bar", Want: false, Name: "unix relative parent"},
		{URL: "~/foo/bar", Want: false, Name: "unix home directory"},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, got := IsValidURL(tt.URL)
			if got != tt.Want {
				t.Errorf("got: %v, want: %v", got, tt.Want)
			}
		})
	}
}

func TestParseURLList(t *testing.T) {
	tests := []struct {
		Name      string
		Input     string
		Filename  string
		Want      []string
		HasErrors bool
		WantDiags hcl.Diagnostics
	}{
		{
			Name:      "empty",
			Input:     "",
			Want:      []string{},
			HasErrors: false,
			WantDiags: nil,
		},
		{
			Name: "all comments",
			Input: `// https://example.com
# https://example.com
;https://example.com

# this is ignored as well
`,
			Want:      []string{},
			HasErrors: false,
			WantDiags: nil,
		},
		{
			Name: "only one",
			Input: `// https://example.com
https://example.com
;https://example.com

# this is ignored as well
`,
			Want: []string{
				"https://example.com",
			},
			HasErrors: false,
			WantDiags: nil,
		},
		{
			Name:  "removes fragment",
			Input: `https://example.com/foo#bar`,
			Want: []string{
				"https://example.com/foo",
			},
			HasErrors: false,
			WantDiags: nil,
		},
		{
			Name:      "invalid url",
			Input:     `not-an-url lol`,
			Want:      nil,
			HasErrors: true,
			WantDiags: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid URL",
					Detail:   "The string 'not-an-url lol' is not a valid url.",
					Subject: &hcl.Range{
						Filename: "list.txt",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: len("not-an-url lol") + 1},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			urls, diags := ParseURLList(test.Input, "list.txt")

			if diags.HasErrors() != test.HasErrors {
				t.Errorf("got: %v, want: %v", diags.HasErrors(), test.HasErrors)
			}

			if !diags.HasErrors() && !reflect.DeepEqual(urls, test.Want) {
				t.Errorf("got: %v, want: %v", urls, test.Want)
			}

			if !diags.HasErrors() && !reflect.DeepEqual(diags, test.WantDiags) {
				t.Errorf("got: %v, want: %v", diags, test.WantDiags)
			}
		})
	}
}

func TestGetURLsFromArgs(t *testing.T) {
	initialWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(initialWd)
	}()

	root := tu.GetOSRoot()
	Fs, Io, Wd = tu.SetupMemMapFs(root)

	Fs.MkdirAll(filepath.Join(root, "other", "directory"), os.ModePerm)
	Fs.MkdirAll(filepath.Join(root, "tmp"), os.ModePerm)
	Io.WriteFile(Fs, filepath.Join(root, "restricted__r.txt"), []byte("not readable"), os.ModePerm)
	Io.WriteFile(Fs, filepath.Join(root, "tmp", "list.ini"), []byte(`// https://example.com
https://example.com
https://more.com?foo=bar#baz
;https://example.com

# this is ignored as well
`), os.ModePerm)
	Io.WriteFile(Fs, filepath.Join(root, "tmp", "invalid.ini"), []byte(`// https://example.com
\x000
;https://example.com

# this is ignored as well
`), os.ModePerm)

	tests := []struct {
		Name      string
		Args      []string
		Want      []string
		WantDiags hcl.Diagnostics
	}{
		{
			Name:      "empty",
			Args:      []string{},
			Want:      []string{},
			WantDiags: nil,
		},
		{
			Name: "file with comments",
			Args: []string{filepath.Join(root, "tmp", "list.ini")},
			Want: []string{
				"https://example.com",
				"https://more.com?foo=bar",
			},
			WantDiags: nil,
		},
		{
			Name: "relative file with comments",
			Args: []string{filepath.Join("tmp", "list.ini")},
			Want: []string{
				"https://example.com",
				"https://more.com?foo=bar",
			},
			WantDiags: nil,
		},
		{
			Name: "one url",
			Args: []string{"https://example.com"},
			Want: []string{
				"https://example.com",
			},
			WantDiags: nil,
		},
		{
			Name: "multiple urls",
			Args: []string{"https://example.com", "http://aws.com"},
			Want: []string{
				"https://example.com",
				"http://aws.com",
			},
			WantDiags: nil,
		},
		{
			Name: "invalid url",
			Args: []string{"no-scheme.org", "https://aws.com"},
			Want: []string{},
			WantDiags: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid argument",
					Detail:   "The argument 'no-scheme.org' is not a valid url, nor a file.",
				},
			},
		},
		{
			Name: "file does not exist",
			Args: []string{"file_does_not_exist.txt"},
			Want: []string{},
			WantDiags: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid argument",
					Detail:   "The argument 'file_does_not_exist.txt' is not a valid url, nor a file.",
				},
			},
		},
		{
			Name: "file not readable",
			Args: []string{"restricted__r.txt"},
			Want: []string{},
			WantDiags: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Could not read file",
					Detail:   "Could not read file '" + filepath.Join(root, "restricted__r.txt") + "'.",
				},
			},
		},
		{
			Name: "url and relative file",
			Args: []string{"https://aws.com", filepath.Join("tmp", "list.ini")},
			Want: []string{
				"https://aws.com",
				"https://example.com",
				"https://more.com?foo=bar",
			},
			WantDiags: nil,
		},
		{
			Name: "url and invalid url in file",
			Args: []string{"https://aws.com", filepath.Join("tmp", "invalid.ini")},
			Want: []string{},
			WantDiags: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid URL",
					Detail:   "The string '\\x000' is not a valid url.",
					Subject: &hcl.Range{
						Filename: filepath.Join(root, "tmp", "invalid.ini"),
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: len("\\x000") + 1},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(tc *testing.T) {
			Wd = filepath.Join(root)

			got, diags := GetURLsFromArgs(test.Args)

			if !slices.Equal(got, test.Want) {
				tc.Errorf("got: %v, want: %v", got, test.Want)
			}

			if !reflect.DeepEqual(diags, test.WantDiags) {
				tc.Errorf("got: %v, want: %v", diags, test.WantDiags)
			}
		})
	}
}
