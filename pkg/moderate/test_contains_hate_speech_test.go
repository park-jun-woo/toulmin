//ff:func feature=moderate type=rule control=iteration dimension=1
//ff:what TestContainsHateSpeech — tests ContainsHateSpeech rule
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestContainsHateSpeech(t *testing.T) {
	tests := []struct {
		name  string
		score float64
		want  bool
	}{
		{"hate", 0.95, true},
		{"not hate", 0.3, false},
		{"borderline", 0.8, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := &ClassifierBacking{Classifier: &mockClassifier{scores: map[string]float64{"hate_speech": tt.score}}}
			ctx := toulmin.NewContext()
			ctx.Set("body", "test")
			got, _ := ContainsHateSpeech(ctx, cb)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
