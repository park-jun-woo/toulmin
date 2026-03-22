//ff:func feature=feature type=util control=sequence
//ff:what TestHashPercentageDeterministic — tests hash percentage is deterministic
package feature

import "testing"

func TestHashPercentageDeterministic(t *testing.T) {
	a := hashPercentage("user-abc")
	b := hashPercentage("user-abc")
	if a != b {
		t.Errorf("expected deterministic: %f != %f", a, b)
	}
}
