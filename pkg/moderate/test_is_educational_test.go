//ff:func feature=moderate type=rule control=sequence
//ff:what TestIsEducational — tests IsEducational rule
package moderate

import "testing"

func TestIsEducational(t *testing.T) {
	ctx := &ContentContext{Channel: &Channel{Type: "education"}}
	got, _ := IsEducational(nil, ctx, nil)
	if !got {
		t.Error("expected true for education channel")
	}
}
