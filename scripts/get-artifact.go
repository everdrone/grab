//go:build exclude

package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/everdrone/grab/cmd"
	"github.com/everdrone/grab/internal/config"
)

func main() {
	OS := os.Getenv("TARGET_GOOS")
	ARCH := os.Getenv("TARGET_GOARCH")

	if OS == "" {
		OS = runtime.GOOS
	}

	if ARCH == "" {
		ARCH = runtime.GOARCH
	}

	// appends an extension on windows
	extension := ""
	if OS == "windows" {
		extension = ".exe"
	}

	fmt.Printf("%s-%s_%s-%s%s", cmd.RootCmd.Name(), config.Version, OS, ARCH, extension)
}
