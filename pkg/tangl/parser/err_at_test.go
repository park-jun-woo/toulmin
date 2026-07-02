//ff:func feature=tangl type=parser control=sequence
//ff:what TestErrAt — tests errAt formats path, line, and message into a single error
package parser

import "testing"

func TestErrAt(t *testing.T) {
	err := errAt("test.md", 42, "unexpected %q", "token")
	want := `test.md:42: unexpected "token"`
	if err == nil || err.Error() != want {
		t.Fatalf("expected %q, got %v", want, err)
	}
}
