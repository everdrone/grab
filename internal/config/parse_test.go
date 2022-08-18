package config

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/everdrone/grab/internal/context"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

func TestValidateSpec(t *testing.T) {
	tests := []struct {
		Name      string
		Input     string
		HasErrors bool
		NumDiags  int
	}{
		{
			Name: "missing global",
			Input: `
site "mysite" {
	test = "mypattern"
	asset "myasset" {
		pattern = "mypattern"
		capture = 1
	}
}`,
			HasErrors: true,
			NumDiags:  1,
		},
		{
			Name: "missing sites",
			Input: `
global {
	location = "mylocation"
}`,
			HasErrors: true,
			NumDiags:  1,
		},
		{
			// This is a valid spec but we check for this in the ValidateConfig() function
			Name: "missing assets",
			Input: `
global {
	location = "mylocation"
}

site "unsplash" {
	test = "mypattern"
}`,
			HasErrors: false,
			NumDiags:  0,
		},
		{
			// This is a valid spec but we check for this in the ValidateConfig() function
			Name: "missing test pattern",
			Input: `
global {
	location = "mylocation"
}

site "unsplash" {
}`,
			HasErrors: true,
			NumDiags:  1,
		},
		{
			// This is a valid spec but we check for this in the ValidateConfig() function
			Name: "missing asset pattern",
			Input: `
global {
	location = "mylocation"
}

site "unsplash" {
	test = "mypattern"

	asset "myasset" {
		capture = 1
	}
}`,
			HasErrors: true,
			NumDiags:  1,
		},
		{
			// This is a valid spec but we check for this in the ValidateConfig() function
			Name: "missing capture and pattern",
			Input: `
global {
	location = "mylocation"
}

site "unsplash" {
	test = "mypattern"

	asset "myasset" {
		pattern = "mypattern"
	}

	asset "myasset2" {
		capture = 1
	}
}`,
			HasErrors: true,
			NumDiags:  2,
		},
		{
			// This is a valid spec but we check for this in the ValidateConfig() function
			Name: "too many transforms",
			Input: `
global {
	location = "mylocation"
}

site "unsplash" {
	test = "mypattern"

	asset "myasset" {
		pattern = "mypattern"
		capture = 1

		transform "mytransform" {
			pattern = "mypattern"
			replace = "myreplace"
		}

		transform "mytransform2" {
			pattern = "mypattern"
			replace = "myreplace"
		}

		transform "mytransform3" {
			pattern = "mypattern"
			replace = "myreplace"
		}
	}
}`,
			HasErrors: true,
			NumDiags:  1,
		},
		{
			// This is a valid spec but we check for this in the ValidateConfig() function
			Name: "invalid attribute",
			Input: `
global {
	location = "mylocation"
}

site "unsplash" {
	test = "mypattern"

	asset "myasset" {
		pattern = "mypattern"
		capture = 1

		something = 2

		transform "mytransform" {
			pattern = "mypattern"
			replace = "myreplace"
		}
	}
}`,
			HasErrors: true,
			NumDiags:  1,
		},
	}

	ctx := context.BuildInitialContext()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			p := hclparse.NewParser()
			file, diags := p.ParseHCL([]byte(test.Input), "test.hcl")
			if diags.HasErrors() {
				t.Fatal(diags)
			}

			diags = ValidateSpec(&file.Body, ctx)
			if diags.HasErrors() != test.HasErrors {
				t.Errorf("got: %v, want: %v", diags.HasErrors(), test.HasErrors)
			}

			if len(diags) != test.NumDiags {
				t.Errorf("expected %d diagnostics, got %d", test.NumDiags, len(diags))
				t.Log(diags)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		Name      string
		Input     string
		HasErrors bool
		WantDiags hcl.Diagnostics
	}{
		{
			Name: "at least one asset or info per site",
			Input: `
site "mysite" {
	test = "mypattern"
}`,
			HasErrors: true,
			WantDiags: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Insufficient \"site\" and \"info\" blocks",
					Detail:   "At least one asset or one info block must be defined inside a \"site\" block.",
				},
			},
		},
		{
			Name: "ok at least one asset",
			Input: `
site "mysite" {
	test = "mypattern"

	asset "myasset" {
		pattern = "x"
	}
}`,
			HasErrors: false,
			WantDiags: nil,
		},
		{
			Name: "ok at least one info",
			Input: `
site "mysite" {
	test = "mypattern"

	info "data" {
		pattern = "x"
	}
}`,
			HasErrors: false,
			WantDiags: nil,
		},
		{
			Name: "invalid transform block label",
			Input: `
site "mysite" {
	test = "mypattern"

	asset "myasset" {
		pattern = "x"
		transform "invalid_here" {
		}
	}
}`,
			HasErrors: true,
			WantDiags: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid block label",
					Detail:   "\"transform\" block labels must be either \"url\" or \"filename\".",
				},
			},
		},
		{
			Name: "duplicate transform block label",
			Input: `
site "mysite" {
	test = "mypattern"

	asset "myasset" {
		pattern = "x"
		transform "url" {
		}
		transform "url" {
		}
		transform "filename" {
		}
	}
}`,
			HasErrors: true,
			WantDiags: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Duplicate block label",
					Detail:   "No more than one \"transform\" block with the label \"url\" is allowed.",
				},
			},
		},
		{
			Name: "ok transform blocks valid",
			Input: `
site "mysite" {
	test = "mypattern"

	asset "myasset" {
		pattern = "x"
		transform "url" {
		}
		transform "filename" {
		}
	}
}`,
			HasErrors: false,
			WantDiags: nil,
		},
		{
			Name: "invalid subdirectory from attribute",
			Input: `
site "mysite" {
	test = "mypattern"

	asset "myasset" {
		pattern = "x"
	}

	subdirectory {
		pattern = "x"
		from = "xxx"
	}
}`,
			HasErrors: true,
			WantDiags: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid block attribute",
					Detail:   "The \"from\" attribute must be either \"body\" or \"url\".",
				},
			},
		},
		{
			Name: "ok subdirectory block valid",
			Input: `
site "mysite" {
	test = "mypattern"

	asset "myasset" {
		pattern = "x"
	}

	subdirectory {
		pattern = "x"
		from = "body"
	}
}`,
			HasErrors: false,
			WantDiags: nil,
		},
	}

	ctx := context.BuildInitialContext()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			p := hclparse.NewParser()
			file, diags := p.ParseHCL([]byte(test.Input), "test.hcl")
			if diags.HasErrors() {
				t.Fatal(diags)
			}

			root := file.Body.(*hclsyntax.Body)

			diags = ValidateConfig(root, ctx)
			if diags.HasErrors() != test.HasErrors {
				t.Fatalf("got: %v, want: %v", diags.HasErrors(), test.HasErrors)
			}

			for i, wantDiag := range test.WantDiags {
				if diags[i].Severity != wantDiag.Severity {
					t.Errorf("got: %v, want: %v", diags[i].Severity, wantDiag.Severity)
				}

				if diags[i].Summary != wantDiag.Summary {
					t.Errorf("got: %v, want: %v", diags[i].Summary, wantDiag.Summary)
				}

				if diags[i].Detail != wantDiag.Detail {
					t.Errorf("got: %v, want: %v", diags[i].Detail, wantDiag.Detail)
				}
			}
		})
	}
}

func TestEvaluateRegexPattern(t *testing.T) {
	tests := []struct {
		Name      string
		Attribute *hclsyntax.Attribute
		WantStr   string
		WantRegex *regexp.Regexp
		HasError  bool
		WantDiags hcl.Diagnostics
	}{
		{
			Name: "ok valid regex",
			Attribute: &hclsyntax.Attribute{
				Name: "pattern",
				Expr: &hclsyntax.LiteralValueExpr{
					Val: cty.StringVal("^/foo/bar$"),
				},
				EqualsRange: hcl.Range{
					Filename: "test.hcl",
					Start:    hcl.Pos{Line: 1, Column: 1, Byte: 0},
					End:      hcl.Pos{Line: 1, Column: 12, Byte: 11},
				},
			},
			WantStr:   "^/foo/bar$",
			WantRegex: regexp.MustCompile("^/foo/bar$"),
			HasError:  false,
			WantDiags: nil,
		},
		{
			Name: "invalid regex",
			Attribute: &hclsyntax.Attribute{
				Name: "pattern",
				Expr: &hclsyntax.LiteralValueExpr{
					Val: cty.StringVal("(?!x)"),
				},
			},
			WantStr:   "",
			WantRegex: nil,
			HasError:  true,
			WantDiags: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid regex pattern",
					Detail:   "error parsing regexp: invalid or unsupported Perl syntax: `(?!`",
				},
			},
		},
	}

	ctx := context.BuildInitialContext()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gotStr, gotRegex, diags := EvaluateRegexPattern(test.Attribute, ctx)
			if diags.HasErrors() != test.HasError {
				t.Fatalf("got: %v, want: %v", diags.HasErrors(), test.HasError)
			}

			if gotStr != test.WantStr {
				t.Errorf("got: %v, want: %v", gotStr, test.WantStr)
			}

			if !reflect.DeepEqual(gotRegex, test.WantRegex) {
				t.Errorf("got: %v, want: %v", gotRegex, test.WantRegex)
			}

			for i, wantDiag := range test.WantDiags {
				if diags[i].Severity != wantDiag.Severity {
					t.Errorf("got: %v, want: %v", diags[i].Severity, wantDiag.Severity)
				}

				if diags[i].Summary != wantDiag.Summary {
					t.Errorf("got: %v, want: %v", diags[i].Summary, wantDiag.Summary)
				}

				if diags[i].Detail != wantDiag.Detail {
					t.Errorf("got: %v, want: %v", diags[i].Detail, wantDiag.Detail)
				}
			}
		})
	}
}

func TestBuildRegexCache(t *testing.T) {
	tests := []struct {
		Name      string
		Input     string
		Want      RegexCacheMap
		WantDiags hcl.Diagnostics
	}{
		{
			Name: "ok single regex",
			Input: `
site "foo" {
	test = "^/foo/bar$"
}`,
			Want: RegexCacheMap{
				"^/foo/bar$": regexp.MustCompile("^/foo/bar$"),
			},
			WantDiags: nil,
		},
		{
			Name: "test look around",
			Input: `
site "foo" {
	test = "(?!look)"
}`,
			Want: RegexCacheMap(nil),
			WantDiags: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid regex pattern",
					Detail:   "error parsing regexp: invalid or unsupported Perl syntax: `(?!`",
				},
			},
		},
		{
			Name: "pattern look around",
			Input: `
site "foo" {
	test = "^abc$"

	asset "bar" {
		pattern = "(?!abc)"
	}
}`,
			Want: RegexCacheMap(nil),
			WantDiags: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid regex pattern",
					Detail:   "error parsing regexp: invalid or unsupported Perl syntax: `(?!`",
				},
			},
		},
		{
			Name: "ok all blocks ok",
			Input: `
site "foo" {
	test = "^abc$"

	asset "bar" {
		pattern = "^abc$"
	}

	info "baz" {
		pattern = "patt(ern)?"
	}

	subdirectory "qux" {
		pattern = "^/qux$"
	}
}`,
			Want: RegexCacheMap{
				"^abc$":      regexp.MustCompile("^abc$"),
				"^/qux$":     regexp.MustCompile("^/qux$"),
				"patt(ern)?": regexp.MustCompile("patt(ern)?"),
			},
			WantDiags: nil,
		},
		{
			Name: "ok all blocks + transforms ok",
			Input: `
site "foo" {
	test = "^abc$"

	asset "bar" {
		pattern = "^abc$"

		transform url {
			pattern = "^def$"
		}

		transform filename {
			pattern = "^ghi$"
		}
	}

	info "baz" {
		pattern = "patt(ern)?"
	}

	subdirectory "qux" {
		pattern = "^/qux$"
	}
}`,
			Want: RegexCacheMap{
				"^abc$":      regexp.MustCompile("^abc$"),
				"^def$":      regexp.MustCompile("^def$"),
				"^ghi$":      regexp.MustCompile("^ghi$"),
				"^/qux$":     regexp.MustCompile("^/qux$"),
				"patt(ern)?": regexp.MustCompile("patt(ern)?"),
			},
			WantDiags: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			p := hclparse.NewParser()
			file, diags := p.ParseHCL([]byte(test.Input), "test.hcl")
			if diags.HasErrors() {
				t.Fatal(diags)
			}

			ctx := context.BuildInitialContext()
			root := file.Body.(*hclsyntax.Body)

			got, diags := BuildRegexCache(root, ctx)
			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf("got: %v, want: %v", got, test.Want)
			}

			for i, wantDiag := range test.WantDiags {
				if diags[i].Severity != wantDiag.Severity {
					t.Errorf("got: %v, want: %v", diags[i].Severity, wantDiag.Severity)
				}

				if diags[i].Summary != wantDiag.Summary {
					t.Errorf("got: %v, want: %v", diags[i].Summary, wantDiag.Summary)
				}

				if diags[i].Detail != wantDiag.Detail {
					t.Errorf("got: %v, want: %v", diags[i].Detail, wantDiag.Detail)
				}
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		Name           string
		Input          string
		Want           *Config
		WantContext    *hcl.EvalContext
		WantRegexCache RegexCacheMap
		WantError      bool
	}{
		{
			Name: "wrong syntax",
			Input: `
site foo" {
	test = "^/foo/bar$"
}`,
			Want:           nil,
			WantContext:    nil,
			WantRegexCache: nil,
			WantError:      true,
		},
		{
			Name: "wrong spec",
			Input: `
site "foo" {
	test = "^/foo/bar$"
}`,
			Want:           nil,
			WantContext:    nil,
			WantRegexCache: nil,
			WantError:      true,
		},
		{
			Name: "wrong static validation",
			Input: `
global {
	location = "some location"
}

site "foo" {
	test = "^/foo/bar$"
}`,
			Want:           nil,
			WantContext:    nil,
			WantRegexCache: nil,
			WantError:      true,
		},
		{
			Name: "wrong regexp syntax",
			Input: `
global {
	location = "some location"
}

site "foo" {
	test = "x"

	info "bar" {
		pattern = "(?<name>baz)"
		capture = 0
	}
}`,
			Want:           nil,
			WantContext:    nil,
			WantRegexCache: nil,
			WantError:      true,
		},
		{
			Name: "ok",
			Input: `
global {
	location = "some location"
}

site "foo" {
	test = "a(x)b"

	asset "bar" {
		pattern = "baz"
		capture = 0
	}

	info "baz" {
		pattern = "patt(ern)?"
		capture = 1
	}
}`,
			Want: &Config{
				Global: GlobalConfig{
					Location: "some location",
				},
				Sites: []SiteConfig{
					{
						Name: "foo",
						Test: "a(x)b",
						Assets: []AssetConfig{
							{
								Name:    "bar",
								Pattern: "baz",
								Capture: "0",
							},
						},
						Infos: []InfoConfig{
							{
								Name:    "baz",
								Pattern: "patt(ern)?",
								Capture: "1",
							},
						},
					},
				},
			},
			WantRegexCache: RegexCacheMap{
				"a(x)b":      regexp.MustCompile("a(x)b"),
				"baz":        regexp.MustCompile("baz"),
				"patt(ern)?": regexp.MustCompile("patt(ern)?"),
			},
			WantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(tc *testing.T) {
			got, _, gotRegexCache, gotDiags := Parse([]byte(tt.Input), "test.hcl")

			if (gotDiags != nil) != tt.WantError {
				tc.Errorf("got: %v, want: %v", gotDiags, tt.WantError)
			}

			if !reflect.DeepEqual(got, tt.Want) {
				tc.Errorf("got: %v, want: %v", got, tt.Want)
			}

			if !reflect.DeepEqual(gotRegexCache, tt.WantRegexCache) {
				tc.Errorf("got: %v, want: %v", gotRegexCache, tt.WantRegexCache)
			}
		})
	}
}
