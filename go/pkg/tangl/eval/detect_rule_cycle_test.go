//ff:func feature=tangl type=engine control=sequence
//ff:what TestDetectRuleCycle — tests acyclic acceptance, shared deps, and cycle reporting
package eval

import (
	"strings"
	"testing"
)

func TestDetectRuleCycle(t *testing.T) {
	lineOf := func(string) int { return 4 }

	acyclic := map[string][]string{"a": {"c"}, "b": {"c"}, "c": nil}
	if err := detectRuleCycle("doc.md", acyclic, lineOf); err != nil {
		t.Fatalf("acyclic graph with a shared dep should pass: %v", err)
	}
	cyclic := map[string][]string{"a": {"b"}, "b": {"a"}}
	err := detectRuleCycle("doc.md", cyclic, lineOf)
	if err == nil {
		t.Fatal("want cycle error")
	}
	if !strings.Contains(err.Error(), "doc.md:4") || !strings.Contains(err.Error(), "rule cycle") {
		t.Fatalf("cycle error should name the path and line: %v", err)
	}
}
