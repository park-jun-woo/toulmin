//ff:func feature=moderate type=engine control=sequence
//ff:what TestMinPostsSpec_SpecName — tests SpecName returns the fixed spec name
package moderate

import "testing"

func TestMinPostsSpec_SpecName(t *testing.T) {
	spec := &MinPostsSpec{MinPosts: 10}
	if got := spec.SpecName(); got != "MinPostsSpec" {
		t.Fatalf("SpecName() = %q, want %q", got, "MinPostsSpec")
	}
}
