//ff:func feature=engine type=engine control=sequence
//ff:what TestStringListSpecCopiesItems — StringListSpec copies its items in and out, and reports its name and validity
package toulmin

import "testing"

func TestStringListSpecCopiesItems(t *testing.T) {
	in := []string{"a"}
	s := NewStringListSpec("list", in)
	in[0] = "changed"
	out := s.StringList()
	out[0] = "also changed"
	if s.StringList()[0] != "a" || s.SpecName() != "list" || s.Validate() != nil {
		t.Fatalf("spec mutated: %#v", s)
	}
}
