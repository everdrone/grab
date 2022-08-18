package utils

import "testing"

func TestFormatMap(t *testing.T) {
	tests := []struct {
		Name       string
		Separator  string
		AlignRight bool
		Map        map[string]string
		Want       string
	}{
		{
			Name:       "empty map",
			Separator:  ":",
			AlignRight: false,
			Map:        map[string]string{},
			Want:       "",
		},
		{
			Name:       "one element",
			Separator:  ":",
			AlignRight: false,
			Map: map[string]string{
				"foo": "bar",
			},
			Want: "foo:bar\n",
		},
		{
			Name:       "two elements",
			Separator:  ":",
			AlignRight: false,
			Map: map[string]string{
				"foooooooo": "bar",
				"baz":       "qux",
			},
			Want: "baz      :qux\nfoooooooo:bar\n",
		},
		{
			Name:       "align right",
			Separator:  ":",
			AlignRight: true,
			Map: map[string]string{
				"foooooooo": "bar",
				"baz":       "qux",
			},
			Want: "      baz:qux\nfoooooooo:bar\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(tc *testing.T) {
			got := FormatMap(tt.Map, tt.Separator, tt.AlignRight)
			if got != tt.Want {
				tc.Errorf("got: %q, want %q", got, tt.Want)
			}
		})
	}
}
