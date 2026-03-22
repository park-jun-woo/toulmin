//ff:func feature=moderate type=rule control=sequence
//ff:what TestContainsNSFW — tests ContainsNSFW rule
package moderate

import "testing"

func TestContainsNSFW(t *testing.T) {
	c := &mockClassifier{scores: map[string]float64{"nsfw": 0.9}}
	got, _ := ContainsNSFW(&Content{Body: "test"}, nil, c)
	if !got {
		t.Error("expected nsfw detected")
	}
}
