package cmd

import (
	"os"
	"testing"

	"github.com/everdrone/grab/internal/utils"
	tu "github.com/everdrone/grab/testutils"
)

func TestGetCmd(t *testing.T) {
	initialWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(initialWd)
	}()

	root := tu.GetOSRoot()
	utils.Fs, utils.AFS, utils.Wd = tu.SetupMemMapFs(root)

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
