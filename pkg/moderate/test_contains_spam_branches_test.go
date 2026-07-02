//ff:func feature=moderate type=rule control=sequence
//ff:what TestContainsSpam_Branches — covers empty specs, non-string body, and non-spam branches of ContainsSpam
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestContainsSpam_Branches(t *testing.T) {
	t.Run("EmptySpecs", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("body", "test")
		got, _ := ContainsSpam(ctx, toulmin.Specs{})
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})

	t.Run("NonStringBody", func(t *testing.T) {
		cb := &ClassifierSpec{Classifier: &mockClassifier{scores: map[string]float64{"spam": 0.9}}}
		ctx := toulmin.NewContext()
		ctx.Set("body", 123)
		got, _ := ContainsSpam(ctx, toulmin.Specs{cb})
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})

	t.Run("BelowThreshold", func(t *testing.T) {
		cb := &ClassifierSpec{Classifier: &mockClassifier{scores: map[string]float64{"spam": 0.2}}}
		ctx := toulmin.NewContext()
		ctx.Set("body", "test")
		got, _ := ContainsSpam(ctx, toulmin.Specs{cb})
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})
}
