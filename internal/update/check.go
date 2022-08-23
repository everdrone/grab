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

	if latest, ok := tagName.(string); ok {
		if latest[0] != 'v' {
			latest = "v" + latest
		}

		if !semver.IsValid(latest) {
			return "", fmt.Errorf("invalid version: %s", latest)
		}

		current := "v" + config.Version

		if semver.Compare(latest, current) == 1 {
			return latest[1:], nil
		}
	} else {
		return "", fmt.Errorf("invalid tag name")
	}

	return "", nil
}
