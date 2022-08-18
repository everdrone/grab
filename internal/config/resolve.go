package config

import (
	"fmt"
	"path/filepath"

	"github.com/everdrone/grab/internal/utils"
)

func getAncestors(path string) []string {
	ancestors := []string{
		path,
	}

	for {
		parent := filepath.Dir(path)
		if parent == path {
			break
		}

		ancestors = append(ancestors, parent)
		path = parent
	}

	return ancestors
}

// NOTE: the path must be absolute!
func Resolve(filename, path string) (string, error) {
	ancestors := getAncestors(path)

	for _, ancestor := range ancestors {
		current := filepath.Join(ancestor, filename)

		exists, err := utils.AFS.Exists(current)
		if err != nil {
			return "", fmt.Errorf("could not check existence of %q: %s", current, err)
		}

		if exists {
			return current, nil
		}
	}

	return "", fmt.Errorf("could not resolve %s", filename)
}
