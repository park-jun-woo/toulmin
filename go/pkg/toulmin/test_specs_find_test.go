//ff:func feature=engine type=engine control=sequence
//ff:what TestSpecsFind — tests Specs.Find returns matching spec and nil for nonexistent
package toulmin

import (
	"testing"
)

func TestSpecsFind(t *testing.T) {
	specs := Specs{&testSpec{Value: "hello"}}

	found := specs.Find("testSpec")
	if found == nil {
		t.Fatal("expected to find testSpec, got nil")
	}
	ts, ok := found.(*testSpec)
	if !ok {
		t.Fatal("expected *testSpec type")
	}
	if ts.Value != "hello" {
		t.Errorf("expected Value 'hello', got '%s'", ts.Value)
	}

	missing := specs.Find("nonexistent")
	if missing != nil {
		t.Errorf("expected nil for nonexistent spec, got %v", missing)
	}
}
