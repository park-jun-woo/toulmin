//ff:func feature=engine type=test control=sequence
//ff:what TestDuplicate — tests duplicate rule registration panics
package toulmin

import (
	"strings"
	"testing"
)

func TestDuplicateRule(t *testing.T) {
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
	g.Rule(WarrantA)
	g.Rule(WarrantA)
}

func TestDuplicateCounter(t *testing.T) {
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
	g.Counter(RebuttalB)
	g.Counter(RebuttalB)
}

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
