//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseDefinitionItemChildren — file and string-list definitions reject child fields; their own errors pass through
package parser

import (
	"strings"
	"testing"
)

func TestParseDefinitionItemChildren(t *testing.T) {
	child := []item{{Text: "`field` means 1", Line: 3}}
	_, err := parseDefinitionItem(item{Text: "`d` reads `JSON list` from \"d.json\"", Line: 2, Children: child}, "t.md")
	if err == nil || !strings.Contains(err.Error(), "file definition cannot have child fields") {
		t.Fatalf("file definition with children: %v", err)
	}
	_, err = parseDefinitionItem(item{Text: "`d` means [\"a\"]", Line: 2, Children: child}, "t.md")
	if err == nil || !strings.Contains(err.Error(), "string-list definition cannot have child fields") {
		t.Fatalf("string-list definition with children: %v", err)
	}
	_, err = parseDefinitionItem(item{Text: "`d` reads `JSON list` \"d.json\"", Line: 2, Children: child}, "t.md")
	if err == nil || !strings.Contains(err.Error(), "expected 'from'") {
		t.Fatalf("file definition parse error must win over the children check: %v", err)
	}
}
