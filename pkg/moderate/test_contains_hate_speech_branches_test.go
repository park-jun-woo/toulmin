//ff:func feature=moderate type=rule control=sequence
//ff:what TestContainsHateSpeech_Branches — covers empty specs and non-string body branches of ContainsHateSpeech
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestContainsHateSpeech_Branches(t *testing.T) {
	t.Run("EmptySpecs", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("body", "test")
		got, _ := ContainsHateSpeech(ctx, toulmin.Specs{})
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})

	t.Run("NonStringBody", func(t *testing.T) {
		cb := &ClassifierSpec{Classifier: &mockClassifier{scores: map[string]float64{"hate_speech": 0.95}}}
		ctx := toulmin.NewContext()
		ctx.Set("body", 123)
		got, _ := ContainsHateSpeech(ctx, toulmin.Specs{cb})
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})
}
