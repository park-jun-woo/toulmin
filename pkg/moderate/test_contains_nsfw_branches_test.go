//ff:func feature=moderate type=rule control=sequence
//ff:what TestContainsNSFW_Branches — covers empty specs, non-string body, and non-nsfw branches of ContainsNSFW
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestContainsNSFW_Branches(t *testing.T) {
	t.Run("EmptySpecs", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("body", "test")
		got, _ := ContainsNSFW(ctx, toulmin.Specs{})
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})

	t.Run("NonStringBody", func(t *testing.T) {
		cb := &ClassifierSpec{Classifier: &mockClassifier{scores: map[string]float64{"nsfw": 0.9}}}
		ctx := toulmin.NewContext()
		ctx.Set("body", 123)
		got, _ := ContainsNSFW(ctx, toulmin.Specs{cb})
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})

	t.Run("BelowThreshold", func(t *testing.T) {
		cb := &ClassifierSpec{Classifier: &mockClassifier{scores: map[string]float64{"nsfw": 0.3}}}
		ctx := toulmin.NewContext()
		ctx.Set("body", "test")
		got, _ := ContainsNSFW(ctx, toulmin.Specs{cb})
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})
}
