//ff:func feature=moderate type=rule control=sequence
//ff:what TestContainsSpam — tests ContainsSpam rule
package moderate

import "testing"

func TestContainsSpam(t *testing.T) {
	cb := &ClassifierBacking{Classifier: &mockClassifier{scores: map[string]float64{"spam": 0.9}}}
	got, _ := ContainsSpam(&Content{Body: "test"}, nil, cb)
	if !got {
		t.Error("expected spam detected")
	}
}
