package context

import (
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

func TestGetEnvironmentMap(t *testing.T) {
	tests := []struct {
		Name   string
		Before func()
		Want   map[string]cty.Value
	}{
		{
			Name: "empty environment",
			Before: func() {
				os.Clearenv()
			},
			Want: map[string]cty.Value{},
		},
		{
			Name: "empty environment",
			Before: func() {
				os.Setenv("foo", "bar")
			},
			Want: map[string]cty.Value{
				"foo": cty.StringVal("bar"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Before()
			got := GetEnvironmentMap()

			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf("got: %v, want: %v", got, test.Want)
			}
		})
	}
}

func TestBuildInitialContext(t *testing.T) {
	tests := []struct {
		Name string
		Want *hcl.EvalContext
	}{
		{
			Name: "body and url constants",
			Want: &hcl.EvalContext{
				Variables: map[string]cty.Value{
					"env": cty.ObjectVal(GetEnvironmentMap()),
					"os": cty.ObjectVal(map[string]cty.Value{
						"name": cty.StringVal(runtime.GOOS),
						"arch": cty.StringVal(runtime.GOARCH),
					}),
					"body": cty.StringVal("body"),
					"url":  cty.StringVal("url"),
				},
				Functions: map[string]function.Function{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got := BuildInitialContext()

			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf("got: %v, want: %v", got, test.Want)
			}
		})
	}
}
