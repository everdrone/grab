package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	changelog, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		log.Fatal("Could not read CHANGELOG.md")
	}

	versionHeader := regexp.MustCompile(`### \[v?\d+\.\d+\.\d+\]`)

	var start int = -1
	var end int = -1

	lines := strings.Split(string(changelog), "\n")
	for i, line := range lines {
		if versionHeader.MatchString(line) {
			if start == -1 {
				start = i + 1
			} else if end == -1 {
				end = i
				break
			}
		}
	}

	body := strings.Join(lines[start:end], "\n")
	body = strings.TrimSpace(body)

	fmt.Fprint(os.Stdout, body)
}
