//ff:func feature=moderate type=engine control=sequence
//ff:what TestReview_SpamTrustedDefeat — tests trusted user defeats spam rebuttal
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestReview_SpamTrustedDefeat(t *testing.T) {
	cb := &ClassifierBacking{Classifier: &mockClassifier{scores: map[string]float64{"spam": 0.9}}}

	g := toulmin.NewGraph("test:trusted-defeat")
	verified := g.Rule(IsVerifiedUser)
	spam := g.Counter(ContainsSpam).Backing(cb)
	trusted := g.Except(IsTrustedUser).Backing(&TrustScoreBacking{MinScore: 0.9})
	spam.Attacks(verified)
	trusted.Attacks(spam)

	mod := NewModerator(g)
	content := &Content{Body: "looks like spam"}
	ctx := &ContentContext{Author: &Author{Verified: true, TrustScore: 0.95}, Channel: &Channel{}}

	result, err := mod.Review(content, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Action == ActionBlock {
		t.Error("expected not block (trusted defeats spam)")
	}
}
