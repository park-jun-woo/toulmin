package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestReview_Allow(t *testing.T) {
	g := toulmin.NewGraph("test:allow")
	g.Warrant(IsVerifiedUser, nil, 1.0)

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

func TestReview_Block(t *testing.T) {
	c := &mockClassifier{scores: map[string]float64{"hate_speech": 0.95}}

	g := toulmin.NewGraph("test:block")
	verified := g.Warrant(IsVerifiedUser, nil, 1.0)
	hate := g.Rebuttal(ContainsHateSpeech, c, 1.0)
	g.Defeat(hate, verified)

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

func TestReview_NewsDefeat(t *testing.T) {
	c := &mockClassifier{scores: map[string]float64{"hate_speech": 0.95}}

	g := toulmin.NewGraph("test:news-defeat")
	verified := g.Warrant(IsVerifiedUser, nil, 1.0)
	hate := g.Rebuttal(ContainsHateSpeech, c, 1.0)
	news := g.Defeater(IsNewsContext, nil, 1.0)
	g.Defeat(hate, verified)
	g.Defeat(news, hate)

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

func TestReview_SpamTrustedDefeat(t *testing.T) {
	c := &mockClassifier{scores: map[string]float64{"spam": 0.9}}

	g := toulmin.NewGraph("test:trusted-defeat")
	verified := g.Warrant(IsVerifiedUser, nil, 1.0)
	spam := g.Rebuttal(ContainsSpam, c, 1.0)
	trusted := g.Defeater(IsTrustedUser, 0.9, 1.0)
	g.Defeat(spam, verified)
	g.Defeat(trusted, spam)

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

func TestReview_TraceIncluded(t *testing.T) {
	g := toulmin.NewGraph("test:trace")
	g.Warrant(IsVerifiedUser, nil, 1.0)

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
