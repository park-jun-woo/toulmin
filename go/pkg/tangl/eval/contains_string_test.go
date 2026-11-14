//ff:func feature=tangl type=engine control=sequence
//ff:what TestContainsString — exact, whitespace-sensitive membership
package eval

import "testing"

func TestContainsString(t *testing.T) {
	items := []string{"a,b", " c "}
	if !containsString(items, "a,b") || !containsString(items, " c ") {
		t.Fatal("exact members must match")
	}
	if containsString(items, "c") || containsString(items, "a") || containsString(nil, "") {
		t.Fatal("partial, trimmed or empty-list lookups must not match")
	}
}
