package utils

import (
	"bytes"
	"fmt"
	"sort"
	"text/tabwriter"
)

func FormatMap(m map[string]string, separator string, alignRight bool) string {
	buf := &bytes.Buffer{}

	flags := uint(0)
	if alignRight {
		flags = flags | tabwriter.AlignRight
	}

	w := tabwriter.NewWriter(buf, 0, 0, 0, ' ', flags)

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for k := range keys {
		fmt.Fprintf(w, "%s\t%s%s\n", keys[k], separator, m[keys[k]])
	}

	w.Flush()
	return buf.String()
}
