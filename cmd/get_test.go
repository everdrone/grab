package cmd

import "testing"

func TestGetCmd(t *testing.T) {
	tests := []struct {
		Name string
	}{
		{
			Name: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
		})
	}
}
