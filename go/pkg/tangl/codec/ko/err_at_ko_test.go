//ff:func feature=tangl type=parser control=sequence
//ff:what TestErrAtKo — formats "path:line: message" with argument substitution
package ko

import "testing"

func TestErrAtKo(t *testing.T) {
	err := errAtKo("t.md", 7, "bad token %q", "x")
	if err == nil {
		t.Fatal("expected a non-nil error")
	}
	if got, want := err.Error(), `t.md:7: bad token "x"`; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
