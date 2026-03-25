//ff:func feature=moderate type=engine control=sequence
//ff:what TestReview_TraceIncluded — tests review result includes trace data
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestReview_TraceIncluded(t *testing.T) {
	g := toulmin.NewGraph("test:trace")
	g.Rule(IsVerifiedUser)

	mod := NewModerator(g)
	content := &Content{Body: "hello"}
	ctx := &ContentContext{Author: &Author{Verified: true}, Channel: &Channel{}}

	result, err := mod.Review(content, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Trace) == 0 {
		t.Error("expected non-empty trace")
	}
}
