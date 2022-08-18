package utils

import (
	"reflect"
	"regexp"
	"testing"
)

func TestHasNamedCaptures(t *testing.T) {
	tests := []struct {
		Name      string
		Regex     string
		Want      bool
		WantNames []string
	}{
		{
			Name:      "no named captures",
			Regex:     "^foo\\s+(bar)?(?:baz)$",
			Want:      false,
			WantNames: []string{"", ""},
		},
		{
			Name:      "named captures",
			Regex:     "^(?P<foo>foo)$",
			Want:      true,
			WantNames: []string{"", "foo"},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			re, err := regexp.Compile(test.Regex)
			if err != nil {
				t.Fatal(err)
			}
			got, gotNames := HasNamedCaptures(re)

			if got != test.Want {
				t.Errorf("got: %v, want: %v", got, test.Want)
			}
			if !reflect.DeepEqual(gotNames, test.WantNames) {
				t.Errorf("got: %v, want: %v", gotNames, test.WantNames)
			}
		})
	}
}

/*
// MARK: - With named captures:
		   with find ALL:

	re := "(?P<first>foo)\s*(?P<second>bar)"
	s := "foo bar foo bar"
	capture := "first"
	matches := [
		["foo bar", "foo", "bar"],
		["foo bar", "foo", "bar"],
	]
	names := ["", "first", "second"]
	output := ["foo", "foo"]

// MARK: - With named captures:
		   with find ONE:

	re := "(?P<first>foo)\s*(?P<second>bar)"
	s := "foo bar foo bar"
	capture := "first"
	matches := ["foo bar", "foo", "bar"]
	names := ["", "first", "second"]
	output := ["foo"]

// MARK: - With capture groups:
		   with find ALL:

	re := "(foo)\s*(bar)"
	s := "foo bar foo bar"
	capture := 1
	matches := [
		["foo bar", "foo", "bar"],
		["foo bar", "foo", "bar"],
	]
	names := ["", "", ""]
	output := ["foo", "foo"]

// MARK: - With capture groups:
		   with find ONE:

	re := "(foo)\s*(bar)"
	s := "foo bar foo bar"
	capture := 1
	matches := ["foo bar", "foo", "bar"]
	names := ["", "", ""]
	output := ["foo"]
*/

func TestGetCaptures(t *testing.T) {
	tests := []struct {
		Name    string
		Regex   string
		FindAll bool
		String  string
		Capture string
		Want    []string
		WantErr string
	}{
		{
			Name:    "all with no named captures",
			Regex:   "(foo\\d)\\s*(bar)",
			FindAll: true,
			Capture: "1",
			String:  "foo1 bar foo2 bar",
			Want:    []string{"foo1", "foo2"},
			WantErr: "",
		},
		{
			Name:    "one no named captures",
			Regex:   "(foo\\d)\\s*(bar)",
			FindAll: false,
			Capture: "1",
			String:  "foo1 bar foo2 bar",
			Want:    []string{"foo1"},
			WantErr: "",
		},
		{
			Name:    "all with named captures",
			Regex:   "(?P<first>foo\\d)\\s*(?P<second>bar)",
			FindAll: true,
			Capture: "first",
			String:  "foo1 bar foo2 bar",
			Want:    []string{"foo1", "foo2"},
			WantErr: "",
		},
		{
			Name:    "one no named captures",
			Regex:   "(?P<first>foo\\d)\\s*(?P<second>bar)",
			FindAll: false,
			Capture: "first",
			String:  "foo1 bar foo2 bar",
			Want:    []string{"foo1"},
			WantErr: "",
		},

		// MARK: - Failing tests

		{
			Name:    "named with number capture input",
			Regex:   "(?P<first>foo\\d)\\s*(?P<second>bar)",
			FindAll: false,
			Capture: "0",
			String:  "foo1 bar foo2 bar",
			Want:    nil,
			WantErr: "capture `0` not found in expression `(?P<first>foo\\d)\\s*(?P<second>bar)`",
		},
		{
			Name:    "named with existing number capture input",
			Regex:   "(?P<first>foo\\d)\\s*(?P<second>bar)",
			FindAll: false,
			Capture: "1",
			String:  "foo1 bar foo2 bar",
			Want:    nil,
			WantErr: "capture `1` not found in expression `(?P<first>foo\\d)\\s*(?P<second>bar)`",
		},
		{
			Name:    "named when using non-named captures",
			Regex:   "(foo\\d)\\s*(bar)",
			FindAll: false,
			Capture: "first",
			String:  "foo1 bar foo2 bar",
			Want:    nil,
			WantErr: "capture group `first` is not a number",
		},
		{
			Name:    "no captures",
			Regex:   "(foo\\d)\\s*(bar)",
			FindAll: false,
			Capture: "1",
			String:  "baz qux",
			Want:    nil,
			WantErr: "no captures found for pattern `(foo\\d)\\s*(bar)`, capture `1`",
		},
		{
			Name:    "no captures with find all",
			Regex:   "(foo\\d)\\s*(bar)",
			FindAll: true,
			Capture: "1",
			String:  "baz qux bez qux",
			Want:    nil,
			WantErr: "no captures found for pattern `(foo\\d)\\s*(bar)`, capture `1`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			re, err := regexp.Compile(tt.Regex)
			if err != nil {
				t.Fatal(err)
			}

			got, err := GetCaptures(re, tt.FindAll, tt.Capture, tt.String)
			if err != nil && tt.WantErr == "" {
				t.Errorf("got: %v, want no error", err)
			}

			if err == nil && tt.WantErr != "" {
				t.Errorf("got no error, want: %v", tt.WantErr)
			}

			if err != nil && tt.WantErr != "" {
				if err.Error() != tt.WantErr {
					t.Errorf("got: %v, want: %v", err, tt.WantErr)
				}
			}

			if !reflect.DeepEqual(got, tt.Want) {
				t.Errorf("got: %#v, want: %#v", got, tt.Want)
			}
		})
	}
}
