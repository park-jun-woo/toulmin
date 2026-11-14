//ff:func feature=tangl type=validator control=sequence
//ff:what TestErrAt — tests errAt formats path:line: message correctly
package validate

import "testing"

func TestErrAt(t *testing.T) {
	err := errAt("doc.md", 12, "case %q: %s missing", "x", "thing")
	want := `doc.md:12: case "x": thing missing`
	if err.Error() != want {
		t.Fatalf("expected %q, got %q", want, err.Error())
	}
}
