//ff:func feature=tangl type=parser control=selection
//ff:what parseRunCheck — classify a Provides/Internal body line as run/check and extract its case name
package en

import (
	"fmt"
	"strings"
)

// parseRunCheck classifies a body item's marker-stripped text: slotKind 1 for a
// `run` line, 2 for a `check` line, isBody false for anything else (e.g. a require
// line handled elsewhere). It mirrors the parser's parseRunCheckItem — the keyword
// must land on a word boundary — and extracts the backtick-quoted case name.
func parseRunCheck(text string) (slotKind int, name string, isBody bool, err error) {
	var rest string
	switch {
	case wordPrefix(text, "run"):
		slotKind, rest = 1, text[len("run"):]
	case wordPrefix(text, "check"):
		slotKind, rest = 2, text[len("check"):]
	default:
		return 0, "", false, nil
	}
	rest = strings.TrimLeft(rest, " \t")
	if !strings.HasPrefix(rest, "`") {
		return 0, "", true, fmt.Errorf("expected backtick-quoted case name")
	}
	end := strings.Index(rest[1:], "`")
	if end < 0 {
		return 0, "", true, fmt.Errorf("unterminated backtick-quoted case name")
	}
	return slotKind, rest[1 : 1+end], true, nil
}
