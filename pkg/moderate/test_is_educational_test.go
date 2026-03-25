//ff:func feature=moderate type=rule control=sequence
//ff:what TestIsEducational — tests IsEducational rule
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsEducational(t *testing.T) {
	ctx := toulmin.NewContext()
	ctx.Set("channel", &Channel{Type: "education"})
	got, _ := IsEducational(ctx, nil)
	if !got {
		t.Error("expected true for education channel")
	}
}
