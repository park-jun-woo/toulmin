//ff:func feature=moderate type=rule control=sequence
//ff:what TestContainsSpam — tests ContainsSpam rule
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestContainsSpam(t *testing.T) {
	cb := &ClassifierSpec{Classifier: &mockClassifier{scores: map[string]float64{"spam": 0.9}}}
	ctx := toulmin.NewContext()
	ctx.Set("body", "test")
	got, _ := ContainsSpam(ctx, toulmin.Specs{cb})
	if !got {
		t.Error("expected spam detected")
	}
}
