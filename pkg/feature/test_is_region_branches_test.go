//ff:func feature=feature type=rule control=sequence
//ff:what TestIsRegion_Branches — covers empty specs and non-string region branches of IsRegion
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsRegion_Branches(t *testing.T) {
	t.Run("EmptySpecs", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("region", "KR")
		got, _ := IsRegion(ctx, toulmin.Specs{})
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})

	t.Run("NonStringRegion", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("region", 123)
		got, _ := IsRegion(ctx, toulmin.Specs{&RegionSpec{Region: "KR"}})
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})
}
