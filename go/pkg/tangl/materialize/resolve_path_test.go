//ff:func feature=tangl type=engine control=sequence
//ff:what TestResolvePath — absolute paths are cleaned, relative ones join the document directory or fail without one
package materialize

import "testing"

func TestResolvePath(t *testing.T) {
	if got, err := resolvePath("", "/data/../lists/x.json"); err != nil || got != "/lists/x.json" {
		t.Fatalf("absolute = %q, %v", got, err)
	}
	if got, err := resolvePath("/docs/sub/p.md", "../x.json"); err != nil || got != "/docs/x.json" {
		t.Fatalf("relative = %q, %v", got, err)
	}
	if _, err := resolvePath("", "x.json"); err == nil || err.Error() != "relative file path requires Document.Path" {
		t.Fatalf("relative without document = %v", err)
	}
}
