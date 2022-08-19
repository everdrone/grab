package testutils

import (
	"strings"
	"testing"
)

func String(v string) *string {
	return &v
}

func Int(v int) *int {
	return &v
}

func Bool(v bool) *bool {
	return &v
}

func AssertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

func EscapeHCLString(s string) string {
	// this does not account for double quote escapes
	// if the string contains \" it will be escaped as \\" and will probably result in an invalid hcl
	// if the string contains \\ it will be escaped as \\\\
	// if the string contains \n it will be escaped as \\n
	return strings.Replace(s, "\\", "\\\\", -1)
}
