package utils

import (
	"reflect"
	"testing"

	tu "github.com/everdrone/grab/testutils"
	"golang.org/x/exp/slices"
)

type testStruct struct {
	Foo string
	Bar int
}

func TestFilter(t *testing.T) {
	tests := []struct {
		Name     string
		Slice    []testStruct
		TestFunc func(testStruct) bool
		Want     []testStruct
	}{
		{
			Name:     "empty slice",
			Slice:    make([]testStruct, 0),
			TestFunc: func(x testStruct) bool { return x.Foo == "foo" },
			Want:     make([]testStruct, 0),
		},
		{
			Name:     "returns empty",
			Slice:    make([]testStruct, 4),
			TestFunc: func(x testStruct) bool { return x.Foo == "something" },
			Want:     make([]testStruct, 0),
		},
		{
			Name:     "returns one",
			Slice:    []testStruct{{"foo", 1}, {"bar", 2}, {"baz", 3}},
			TestFunc: func(x testStruct) bool { return x.Foo == "foo" },
			Want:     []testStruct{{"foo", 1}},
		},
		{
			Name:     "returns many",
			Slice:    []testStruct{{"foo", 1}, {"bar", 2}, {"baz", 3}},
			TestFunc: func(x testStruct) bool { return x.Bar > 1 },
			Want:     []testStruct{{"bar", 2}, {"baz", 3}},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got := Filter(test.Slice, test.TestFunc)

			if !slices.Equal(got, test.Want) {
				t.Errorf("got: %v, want: %v", got, test.Want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		Name  string
		Slice []string
		Value string
		Want  bool
	}{
		{
			Name:  "empty slice",
			Slice: make([]string, 0),
			Value: "foo",
			Want:  false,
		},
		{
			Name:  "contains",
			Slice: []string{"foo", "bar", "baz"},
			Value: "foo",
			Want:  true,
		},
		{
			Name:  "does not contain",
			Slice: []string{"foo", "bar", "baz"},
			Value: "qux",
			Want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got := Contains(tt.Slice, tt.Value)
			if got != tt.Want {
				t.Errorf("got: %v, want: %v", got, tt.Want)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		Name  string
		Slice []string
		Want  []string
	}{
		{
			Name:  "empty slice",
			Slice: make([]string, 0),
			Want:  make([]string, 0),
		},
		{
			Name:  "returns one",
			Slice: []string{"foo", "foo", "foo"},
			Want:  []string{"foo"},
		},
		{
			Name:  "returns many",
			Slice: []string{"foo", "bar", "baz"},
			Want:  []string{"foo", "bar", "baz"},
		},
		{
			Name:  "returns many with duplicates",
			Slice: []string{"foo", "bar", "baz", "bar", "bar", "baz"},
			Want:  []string{"foo", "bar", "baz"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got := Unique(tt.Slice)

			if !slices.Equal(got, tt.Want) {
				t.Errorf("got: %v, want: %v", got, tt.Want)
			}
		})
	}
}

func TestAny(t *testing.T) {
	tests := []struct {
		Name  string
		Slice []string
		Test  func(string) bool
		Want  bool
	}{
		{
			Name:  "empty slice",
			Slice: make([]string, 0),
			Test:  func(x string) bool { return x == "foo" },
			Want:  false,
		},
		{
			Name:  "returns true",
			Slice: []string{"foo", "bar", "baz", "foo"},
			Test:  func(x string) bool { return x == "foo" },
			Want:  true,
		},
		{
			Name:  "returns false",
			Slice: []string{"foo", "bar", "baz", "foo"},
			Test:  func(x string) bool { return x == "qux" },
			Want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got := Any(tt.Slice, tt.Test)
			if got != tt.Want {
				t.Errorf("got: %v, want: %v", got, tt.Want)
			}
		})
	}
}

func TestAll(t *testing.T) {
	tests := []struct {
		Name  string
		Slice []string
		Test  func(string) bool
		Want  bool
	}{
		{
			Name:  "empty slice",
			Slice: make([]string, 0),
			Test:  func(x string) bool { return x == "foo" },
			Want:  false,
		},
		{
			Name:  "some",
			Slice: []string{"foo", "bar", "baz", "foo"},
			Test:  func(x string) bool { return x == "foo" },
			Want:  false,
		},
		{
			Name:  "none",
			Slice: []string{"foo", "bar", "baz", "foo"},
			Test:  func(x string) bool { return x == "qux" },
			Want:  false,
		},
		{
			Name:  "all",
			Slice: []string{"foo", "foo", "foo", "foo"},
			Test:  func(x string) bool { return x == "foo" },
			Want:  true,
		},
		{
			Name:  "one",
			Slice: []string{"foo", "bar", "baz", "foo"},
			Test:  func(x string) bool { return x == "bar" },
			Want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got := All(tt.Slice, tt.Test)
			if got != tt.Want {
				t.Errorf("got: %v, want: %v", got, tt.Want)
			}
		})
	}
}

func TestZipMap(t *testing.T) {
	tests := []struct {
		Name   string
		Keys   []string
		Vals   []string
		Want   map[string]string
		Panics bool
	}{
		{
			Name: "empty",
			Keys: make([]string, 0),
			Vals: make([]string, 0),
			Want: make(map[string]string),
		},
		{
			Name: "one",
			Keys: []string{"foo"},
			Vals: []string{"bar"},
			Want: map[string]string{"foo": "bar"},
		},
		{
			Name: "many",
			Keys: []string{"foo", "bar"},
			Vals: []string{"baz", "qux"},
			Want: map[string]string{"foo": "baz", "bar": "qux"},
		},
		{
			Name:   "longer keys",
			Keys:   []string{"foo", "bar", "baz"},
			Vals:   []string{"qux"},
			Want:   nil,
			Panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Panics {
				tu.AssertPanic(t, func() { ZipMap(tt.Keys, tt.Vals) })
			} else {
				got := ZipMap(tt.Keys, tt.Vals)
				if !reflect.DeepEqual(got, tt.Want) {
					t.Errorf("got: %v, want: %v", got, tt.Want)
				}
			}
		})
	}
}
