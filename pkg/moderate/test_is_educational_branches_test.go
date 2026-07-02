//ff:func feature=moderate type=rule control=sequence
//ff:what TestIsEducational_Branches — covers non-Channel channel and non-education type branches of IsEducational
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsEducational_Branches(t *testing.T) {
	t.Run("NonChannel", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("channel", "not-a-channel")
		got, _ := IsEducational(ctx, nil)
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})

	t.Run("NotEducation", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("channel", &Channel{Type: "general"})
		got, _ := IsEducational(ctx, nil)
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})
}
