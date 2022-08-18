package utils

import (
	"fmt"
	"regexp"
	"strconv"

	"golang.org/x/exp/slices"
)

func HasNamedCaptures(re *regexp.Regexp) (bool, []string) {
	names := re.SubexpNames()
	return Any(names, func(name string) bool { return name != "" }), names
}

// "findAll" gives back multiple matches. if set to false, only the first match is captured
// capture is either the group index as a string or the group name
func GetCaptures(re *regexp.Regexp, findAll bool, capture string, s string) ([]string, error) {
	named, names := HasNamedCaptures(re)

	captureIndex := 0
	if named {
		captureIndex = slices.Index(names, capture)
		if captureIndex == -1 {
			return nil, fmt.Errorf("capture `%s` not found in expression `%s`", capture, re.String())
		}
	} else {
		integer, err := strconv.Atoi(capture)
		if err != nil {
			return nil, fmt.Errorf("capture group `%s` is not a number", capture)
		}

		captureIndex = integer
	}

	result := make([]string, 0)

	if findAll {
		matches := re.FindAllStringSubmatch(s, -1)
		for _, match := range matches {
			if len(match) > captureIndex {
				result = append(result, match[captureIndex])
			}
		}
	} else {
		match := re.FindStringSubmatch(s)
		if len(match) > captureIndex {
			result = append(result, match[captureIndex])
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no captures found for pattern `%s`, capture `%s`", re.String(), capture)
	}

	return result, nil
}
