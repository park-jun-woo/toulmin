//ff:func feature=moderate type=engine control=sequence
//ff:what TestMinPostsSpec_Validate — tests validation of MinPostsSpec non-negative bound
package moderate

import "testing"

func TestMinPostsSpec_Validate(t *testing.T) {
	if err := (&MinPostsSpec{MinPosts: -1}).Validate(); err == nil {
		t.Fatal("expected error for negative MinPosts")
	}
	if err := (&MinPostsSpec{MinPosts: 5}).Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
