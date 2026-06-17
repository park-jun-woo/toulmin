//ff:func feature=engine type=engine control=sequence
//ff:what TestDuplicateExcept — tests duplicate Except registration panics
package toulmin

import (
	"strings"
	"testing"
)

func TestDuplicateExcept(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
		msg, ok := r.(string)
		if !ok || !strings.Contains(msg, "duplicate rule registration") {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	g := NewGraph("test")
	g.Except(DefeaterC)
	g.Except(DefeaterC)
}
