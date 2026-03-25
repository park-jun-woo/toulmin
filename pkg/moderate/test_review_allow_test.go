//ff:func feature=moderate type=engine control=sequence
//ff:what TestReview_Allow — tests moderator allows verified user content
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestReview_Allow(t *testing.T) {
	g := toulmin.NewGraph("test:allow")
	g.Rule(IsVerifiedUser)

	mod := NewModerator(g)
	content := &Content{Body: "hello"}
	ctx := &ContentContext{Author: &Author{Verified: true}, Channel: &Channel{Type: "general"}}

	result, err := mod.Review(content, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Action != ActionAllow {
		t.Errorf("expected allow, got %s", result.Action)
	}
}
