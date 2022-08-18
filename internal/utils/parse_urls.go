package utils

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

func ParseURLList(contents, filename string) ([]string, hcl.Diagnostics) {
	urls := make([]string, 0)

	lines := strings.Split(strings.ReplaceAll(contents, "\r\n", "\n"), "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" ||
			strings.HasPrefix(line, "#") ||
			strings.HasPrefix(line, "//") ||
			strings.HasPrefix(line, ";") {
			// ignore #, ;, // and empty lines
			continue
		}

		// trim again in case the comment starts with a space
		line = strings.TrimSpace(line)

		url, err := url.Parse(line)
		if err != nil || !url.IsAbs() {
			return nil, hcl.Diagnostics{&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid URL",
				Detail:   fmt.Sprintf("The string '%s' is not a valid url.", line),
				Subject: &hcl.Range{
					Filename: filename,
					Start:    hcl.Pos{Line: i + 1, Column: 1},
					End:      hcl.Pos{Line: i + 1, Column: len(line) + 1},
				},
			}}
		}

		// remove fragment
		url.Fragment = ""
		url.RawFragment = ""

		urls = append(urls, url.String())
	}

	return urls, nil
}

func GetURLsFromArgs(args []string) ([]string, hcl.Diagnostics) {
	urls := make([]string, 0)

	for _, arg := range args {
		absolute := Abs(arg)
		exists, err := AFS.Exists(absolute)
		if err != nil || !exists {
			// check url
			url, err := url.Parse(arg)
			if err != nil || !url.IsAbs() {
				return nil, hcl.Diagnostics{&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid argument",
					Detail:   fmt.Sprintf("The argument '%s' is not a valid url, nor a file.", arg),
				}}
			}

			// remove fragment
			url.Fragment = ""
			url.RawFragment = ""

			urls = append(urls, url.String())
		} else {
			// read file
			fc, err := AFS.ReadFile(absolute)
			if err != nil {
				return nil, hcl.Diagnostics{&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid argument",
					Detail:   fmt.Sprintf("The argument %s is not a valid url, nor a file.", arg),
				}}
			}

			parsed, diags := ParseURLList(string(fc), absolute)
			if diags.HasErrors() {
				return nil, diags
			}

			urls = append(urls, parsed...)
		}
	}

	return urls, nil
}
