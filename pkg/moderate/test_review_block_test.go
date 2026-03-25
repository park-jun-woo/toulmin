//ff:func feature=moderate type=engine control=sequence
//ff:what TestReview_Block — tests moderator blocks hate speech content
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestReview_Block(t *testing.T) {
	cb := &ClassifierBacking{Classifier: &mockClassifier{scores: map[string]float64{"hate_speech": 0.95}}}

	g := toulmin.NewGraph("test:block")
	verified := g.Rule(IsVerifiedUser)
	hate := g.Counter(ContainsHateSpeech).Backing(cb)
	hate.Attacks(verified)

	mod := NewModerator(g)
	content := &Content{Body: "hate content"}
	ctx := &ContentContext{Author: &Author{Verified: true}, Channel: &Channel{Type: "general"}}

	result, err := mod.Review(content, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Action != ActionBlock {
		t.Errorf("expected block, got %s", result.Action)
	}
}
