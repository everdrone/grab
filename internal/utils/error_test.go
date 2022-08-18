package utils

import (
	"bytes"
	"testing"

	"github.com/hashicorp/hcl/v2"
)

func TestPlural(t *testing.T) {
	tests := []struct {
		Singular string
		Plural   string
		Count    int
		Want     string
	}{
		{
			Singular: "foo",
			Plural:   "bar",
			Count:    0,
			Want:     "bar",
		},
		{
			Singular: "foo",
			Plural:   "bar",
			Count:    1,
			Want:     "foo",
		},
		{
			Singular: "foo",
			Plural:   "bar",
			Count:    3,
			Want:     "bar",
		},
	}

	for _, tt := range tests {
		if got := Plural(tt.Count, tt.Singular, tt.Plural); got != tt.Want {
			t.Errorf("got: %q, want %q", got, tt.Want)
		}
	}
}

func TestPrintDiag(t *testing.T) {
	tests := []struct {
		Name string
		Diag *hcl.Diagnostic
		Want string
	}{
		{
			Name: "no subject",
			Diag: &hcl.Diagnostic{
				Severity: DiagError,
				Summary:  "foo",
				Detail:   "bar",
			},
			Want: `╷ Error: foo
╵   bar
`,
		},
		{
			Name: "error with subject",
			Diag: &hcl.Diagnostic{
				Severity: DiagError,
				Summary:  "foo",
				Detail:   "bar",
				Subject:  &hcl.Range{Filename: "foo.hcl", Start: hcl.Pos{Line: 1, Column: 1}, End: hcl.Pos{Line: 1, Column: 2}},
			},
			Want: `╷ Error: foo
│   bar
╵   foo.hcl:1,1-2
`,
		},
		{
			Name: "warning",
			Diag: &hcl.Diagnostic{
				Severity: DiagWarning,
				Summary:  "foo",
				Detail:   "bar",
			},
			Want: `╷ Warning: foo
╵   bar
`,
		},
		{
			Name: "invalid",
			Diag: &hcl.Diagnostic{
				Severity: DiagInvalid,
				Summary:  "foo",
				Detail:   "bar",
			},
			Want: `╷ Invalid: foo
╵   bar
`,
		},
		{
			Name: "info",
			Diag: &hcl.Diagnostic{
				Severity: DiagInfo,
				Summary:  "foo",
				Detail:   "bar",
			},
			Want: `╷ Info: foo
╵   bar
`,
		},
		{
			Name: "debug",
			Diag: &hcl.Diagnostic{
				Severity: DiagDebug,
				Summary:  "foo",
				Detail:   "bar",
			},
			Want: `╷ Debug: foo
╵   bar
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var buf bytes.Buffer
			PrintDiag(&buf, tt.Diag)
			if got := buf.String(); got != tt.Want {
				t.Errorf("got: %s, want %s", got, tt.Want)
			}
		})
	}
}
