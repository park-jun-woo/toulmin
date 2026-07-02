//ff:func feature=tangl type=parser control=sequence
//ff:what parseSubjectSection — parse the tangl:Subject section's namespace
package parser

import "strings"

// parseSubjectSection parses "this document is `<namespace>`".
func parseSubjectSection(sec section, path string) (string, error) {
	items, err := parseItems(sec.Lines, sec.LineOffset, path)
	if err != nil {
		return "", err
	}
	if len(items) == 0 {
		return "", errAt(path, sec.HeaderLine, "tangl:Subject requires 'this document is `name`'")
	}
	it := items[0]
	rest, ok := takeKeyword(it.Text, "this document is")
	if !ok {
		return "", errAt(path, it.Line, "expected 'this document is `name`', got %q", it.Text)
	}
	name, rest, ok := takeBacktick(rest)
	if !ok {
		return "", errAt(path, it.Line, "expected backtick-quoted subject name")
	}
	if strings.TrimSpace(rest) != "" {
		return "", errAt(path, it.Line, "unexpected trailing text: %q", rest)
	}
	return name, nil
}
