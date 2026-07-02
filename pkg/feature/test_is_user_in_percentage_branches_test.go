//ff:func feature=feature type=rule control=sequence
//ff:what TestIsUserInPercentage_Branches — covers empty specs and non-string id branches of IsUserInPercentage
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsUserInPercentage_Branches(t *testing.T) {
	t.Run("EmptySpecs", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("id", "user-123")
		got, _ := IsUserInPercentage(ctx, toulmin.Specs{})
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})

	t.Run("NonStringID", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("id", 123)
		got, _ := IsUserInPercentage(ctx, toulmin.Specs{&PercentageSpec{Percentage: 1.0}})
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})
}
