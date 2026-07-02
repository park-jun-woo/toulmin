//ff:func feature=tangl type=parser control=sequence
//ff:what parseRunCheckItem — parse a "run `case`" or "check `case`" statement
package parser

// parseRunCheckItem parses ("run" | "check") NAME, shared by Provides and
// Internal bodies. kind is "run" or "check" when ok is true.
func parseRunCheckItem(it item, path string) (name string, kind string, ok bool, err error) {
	if rest, matched := takeKeyword(it.Text, "run"); matched {
		return parseRunCheckName(rest, "run", it.Line, path)
	}
	if rest, matched := takeKeyword(it.Text, "check"); matched {
		return parseRunCheckName(rest, "check", it.Line, path)
	}
	return "", "", false, nil
}
