package utils

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/hashicorp/hcl/v2"
	"github.com/mattn/go-isatty"
)

var ErrSilent = errors.New("ErrSilent")

const (
	DiagInvalid hcl.DiagnosticSeverity = iota
	DiagError
	DiagWarning
	DiagInfo
	DiagDebug
)

func PrintDiag(w io.Writer, diag *hcl.Diagnostic) {
	if w, ok := w.(*os.File); !ok || !isatty.IsTerminal(w.Fd()) {
		// disable color if we're not in a terminal
		color.NoColor = true
	}

	top := "╷"
	main := "│"
	bottom := "╵"

	colorizer := color.New()
	subjectColorizer := color.New(color.FgHiBlack)

	var severity string
	switch diag.Severity {
	case DiagInvalid:
		severity = "Invalid"
		colorizer.Add(color.FgHiRed)
	case DiagError:
		severity = "Error"
		colorizer.Add(color.FgRed)
	case DiagWarning:
		severity = "Warning"
		colorizer.Add(color.FgYellow)
	case DiagInfo:
		severity = "Info"
		colorizer.Add(color.FgGreen)
	case DiagDebug:
		severity = "Debug"
		colorizer.Add(color.FgCyan)
	}

	severity = colorizer.Sprint(severity)
	top = colorizer.Sprint(top)
	main = colorizer.Sprint(main)
	bottom = colorizer.Sprint(bottom)

	var str string
	if diag.Subject != nil {
		subject := diag.Subject.String()
		subject = subjectColorizer.Sprint(subject)

		str = fmt.Sprintf(`%s %s: %s
%s   %s
%s   %s
`, top, severity, diag.Summary, main, diag.Detail, bottom, subject)
	} else {
		str = fmt.Sprintf(`%s %s: %s
%s   %s
`, top, severity, diag.Summary, bottom, diag.Detail)
	}

	w.Write([]byte(str))
}

func Plural(count int, singular string, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}
