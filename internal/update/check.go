package update

import (
	"encoding/json"
	"fmt"

	"github.com/everdrone/grab/internal/config"
	"github.com/everdrone/grab/internal/net"
	"golang.org/x/mod/semver"
)

func CheckForUpdates() (string, error) {
	resp, err := net.Fetch(config.LatestReleaseURL, &net.FetchOptions{
		Headers: map[string]string{
			"Accept": "application/vnd.github+json",
		},
		Timeout: 1000,
		Retries: 1,
	})

	if err != nil {
		return "", err
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal([]byte(resp), &decoded); err != nil {
		return "", err
	}

	tagName := decoded["tag_name"]
	if tagName == "" {
		return "", fmt.Errorf("no tag name")
	}

	if version, ok := tagName.(string); ok {
		if version[0] != 'v' {
			version = "v" + version
		}

		if !semver.IsValid(version) {
			return "", fmt.Errorf("invalid version: %s", version)
		}

		if semver.Compare(version, "v"+config.Version) == 1 {
			return version, nil
		}
	} else {
		return "", fmt.Errorf("invalid tag name")
	}

	return "", nil
}
