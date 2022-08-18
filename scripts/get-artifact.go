//go:build exclude

package main

import (
	"fmt"
	"os"

	"github.com/everdrone/grab/cmd"
	"github.com/everdrone/grab/internal/config"
)

func main() {
	goOS := os.Getenv("TARGET_GOOS")
	goArch := os.Getenv("TARGET_GOARCH")

	extension := ""
	if goOS == "windows" {
		extension = ".exe"
	}

	fmt.Printf("%s-%s_%s-%s%s", cmd.RootCmd.Name(), config.Version, goOS, goArch, extension)
}
