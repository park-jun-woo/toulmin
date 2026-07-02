//ff:func feature=moderate type=rule control=sequence
//ff:what TestIsNewsContext_NonChannel — covers non-Channel channel branch of IsNewsContext
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsNewsContext_NonChannel(t *testing.T) {
	ctx := toulmin.NewContext()
	ctx.Set("channel", "not-a-channel")
	got, _ := IsNewsContext(ctx, nil)
	if got != false {
		t.Errorf("got %v, want false", got)
	}
}
