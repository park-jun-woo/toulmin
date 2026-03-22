//ff:func feature=moderate type=rule control=sequence
//ff:what TestIsAdultChannel — tests IsAdultChannel rule
package moderate

import "testing"

func TestIsAdultChannel(t *testing.T) {
	ctx := &ContentContext{Channel: &Channel{AgeGated: true}}
	got, _ := IsAdultChannel(nil, ctx, nil)
	if !got {
		t.Error("expected true for age-gated channel")
	}
}
