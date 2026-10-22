//ff:func feature=tangl type=parser control=sequence
//ff:what parseRunCheckName — parse the backtick-quoted case name following "run"/"check"
package parser

import "strings"

// parseRunCheckName parses the backtick-quoted case name following "run"/"check".
func parseRunCheckName(rest string, kind string, line int, path string) (string, string, bool, error) {
	name, rest, ok := takeBacktick(rest)
	if !ok {
		return "", "", true, errAt(path, line, "expected backtick-quoted case name after %q", kind)
	}
	if strings.TrimSpace(rest) != "" {
		return "", "", true, errAt(path, line, "unexpected trailing text: %q", rest)
	}
	return name, kind, true, nil
}
