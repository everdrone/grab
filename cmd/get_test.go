package cmd

import (
	"os"
	"testing"

	"github.com/everdrone/grab/internal/testutils"
	"github.com/everdrone/grab/internal/utils"
)

func TestGetCmd(t *testing.T) {
	initialWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(initialWd)
	}()

	root := testutils.GetOSRoot()
	utils.Fs, utils.AFS, utils.Wd = testutils.SetupMemMapFs(root)

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
