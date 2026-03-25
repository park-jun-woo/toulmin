//ff:func feature=moderate type=rule control=sequence
//ff:what TestContainsNSFW — tests ContainsNSFW rule
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestContainsNSFW(t *testing.T) {
	cb := &ClassifierSpec{Classifier: &mockClassifier{scores: map[string]float64{"nsfw": 0.9}}}
	ctx := toulmin.NewContext()
	ctx.Set("body", "test")
	got, _ := ContainsNSFW(ctx, toulmin.Specs{cb})
	if !got {
		t.Error("expected nsfw detected")
	}
}
