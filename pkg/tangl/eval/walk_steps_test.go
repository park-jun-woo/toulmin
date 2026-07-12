//ff:func feature=tangl type=engine control=sequence
//ff:what TestWalkSteps — tests no-step identity, multi-step walking, and a missing step
package eval

import "testing"

func TestWalkSteps(t *testing.T) {
	if v, ok := walkSteps("as-is", nil); !ok || v != "as-is" {
		t.Fatalf("no steps should return cur unchanged: got %v,%v", v, ok)
	}
	root := map[string]any{"a": map[string]any{"b": 2}}
	if v, ok := walkSteps(root, []string{"a", "b"}); !ok || v != 2 {
		t.Fatalf("two-step walk: got %v,%v", v, ok)
	}
	if _, ok := walkSteps(root, []string{"a", "ghost"}); ok {
		t.Fatal("missing step should yield ok=false")
	}
}
