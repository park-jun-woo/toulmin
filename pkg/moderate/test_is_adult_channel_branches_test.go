//ff:func feature=moderate type=rule control=sequence
//ff:what TestIsAdultChannel_NonChannel — covers non-Channel channel branch of IsAdultChannel
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsAdultChannel_NonChannel(t *testing.T) {
	ctx := toulmin.NewContext()
	ctx.Set("channel", "not-a-channel")
	got, _ := IsAdultChannel(ctx, nil)
	if got != false {
		t.Errorf("got %v, want false", got)
	}
}
