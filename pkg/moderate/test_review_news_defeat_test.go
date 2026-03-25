//ff:func feature=moderate type=engine control=sequence
//ff:what TestReview_NewsDefeat — tests news context defeats hate speech rebuttal
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestReview_NewsDefeat(t *testing.T) {
	cb := &ClassifierSpec{Classifier: &mockClassifier{scores: map[string]float64{"hate_speech": 0.95}}}

	g := toulmin.NewGraph("test:news-defeat")
	verified := g.Rule(IsVerifiedUser)
	hate := g.Counter(ContainsHateSpeech).With(cb)
	news := g.Except(IsNewsContext)
	hate.Attacks(verified)
	news.Attacks(hate)

	mod := NewModerator(g)
	content := &Content{Body: "news quote with hate speech"}
	ctx := &ContentContext{Author: &Author{Verified: true}, Channel: &Channel{Type: "news"}}

	result, err := mod.Review(content, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Action == ActionBlock {
		t.Error("expected not block (news defeats hate)")
	}
}
