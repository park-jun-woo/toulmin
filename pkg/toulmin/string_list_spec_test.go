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
