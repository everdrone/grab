package net

import (
	"reflect"
	"testing"

	"github.com/everdrone/grab/internal/config"
	"github.com/everdrone/grab/internal/testutils"
)

func TestMergeFetchOptionsChain(t *testing.T) {
	tests := []struct {
		Name     string
		Root     *config.RootNetworkConfig
		Children []*config.NetworkConfig
		Want     *FetchOptions
	}{
		{
			Name:     "all nil results in default",
			Root:     nil,
			Children: []*config.NetworkConfig(nil),
			Want: &FetchOptions{
				Timeout: 3000,
				Retries: 1,
				Headers: make(map[string]string, 0),
			},
		},
		{
			Name: "parent nil results child",
			Root: nil,
			Children: []*config.NetworkConfig{
				{
					Timeout: testutils.Int(20000),
					Retries: testutils.Int(2),
					Headers: &map[string]string{
						"foo": "bar",
					},
				},
			},
			Want: &FetchOptions{
				Timeout: 20000,
				Retries: 2,
				Headers: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			Name: "children nil results parent",
			Root: &config.RootNetworkConfig{
				Timeout: testutils.Int(6000),
				Retries: testutils.Int(3),
				Headers: &map[string]string{
					"foo": "bar",
				},
			},
			Children: []*config.NetworkConfig{},
			Want: &FetchOptions{
				Timeout: 6000,
				Retries: 3,
				Headers: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			Name: "both parent and children with all properties set",
			Root: &config.RootNetworkConfig{
				Timeout: testutils.Int(6000),
				Retries: testutils.Int(3),
				Headers: &map[string]string{
					"foo": "bar",
				},
			},
			Children: []*config.NetworkConfig{
				{
					Timeout: testutils.Int(20000),
					Retries: testutils.Int(2),
					Headers: &map[string]string{
						"baz": "qux",
					},
				},
				{
					Timeout: testutils.Int(4000),
					Retries: testutils.Int(5),
					Headers: &map[string]string{
						"bar": "quf",
					},
				},
			},
			Want: &FetchOptions{
				Timeout: 4000,
				Retries: 5,
				Headers: map[string]string{
					"foo": "bar",
					"baz": "qux",
					"bar": "quf",
				},
			},
		},
		{
			Name: "both parent and children with the same property unset",
			Root: &config.RootNetworkConfig{
				Timeout: testutils.Int(6000),
				Retries: testutils.Int(3),
				Headers: nil,
			},
			Children: []*config.NetworkConfig{
				{
					Timeout: testutils.Int(20000),
					Retries: testutils.Int(2),
					Headers: nil,
				},
				{
					Timeout: testutils.Int(4000),
					Retries: testutils.Int(5),
					Headers: nil,
				},
			},
			Want: &FetchOptions{
				Timeout: 4000,
				Retries: 5,
				Headers: make(map[string]string, 0),
			},
		},
		{
			Name: "both parent and children with the some properties unset",
			Root: &config.RootNetworkConfig{
				Retries: testutils.Int(3),
			},
			Children: []*config.NetworkConfig{
				{
					Timeout: testutils.Int(20000),
					Headers: &map[string]string{
						"foo": "bar",
					},
				},
				{
					Retries: testutils.Int(3),
					Headers: &map[string]string{
						"foo": "baz",
					},
				},
			},
			Want: &FetchOptions{
				Timeout: 20000,
				Retries: 3,
				Headers: map[string]string{
					"foo": "baz",
				},
			},
		},
		{
			Name: "does not inherit",
			Root: &config.RootNetworkConfig{
				Timeout: testutils.Int(6000),
				Retries: testutils.Int(3),
				Headers: &map[string]string{
					"foo": "bar",
				},
			},
			Children: []*config.NetworkConfig{
				{
					Timeout: testutils.Int(20000),
					Retries: testutils.Int(2),
					Headers: &map[string]string{
						"baz": "qux",
					},
				},
				{
					Inherit: testutils.Bool(false),
					Retries: testutils.Int(5),
					Headers: &map[string]string{
						"bar": "quf",
					},
				},
			},
			Want: &FetchOptions{
				Timeout: 3000,
				Retries: 5,
				Headers: map[string]string{
					"bar": "quf",
				},
			},
		},
		{
			Name: "does not inherit multiple",
			Root: &config.RootNetworkConfig{
				Timeout: testutils.Int(6000),
				Retries: testutils.Int(3),
				Headers: &map[string]string{
					"foo": "bar",
				},
			},
			Children: []*config.NetworkConfig{
				{
					Inherit: testutils.Bool(false),
					Timeout: testutils.Int(20000),
					Retries: testutils.Int(2),
					Headers: &map[string]string{
						"baz": "qux",
					},
				},
				{
					Retries: testutils.Int(5),
					Headers: &map[string]string{
						"bar": "quf",
					},
				},
			},
			Want: &FetchOptions{
				Timeout: 20000,
				Retries: 5,
				Headers: map[string]string{
					"baz": "qux",
					"bar": "quf",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if got := MergeFetchOptionsChain(tt.Root, tt.Children...); !reflect.DeepEqual(got, tt.Want) {
				t.Errorf("got: %v, want: %v", got, tt.Want)
			}
		})
	}
}
