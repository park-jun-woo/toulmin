//ff:func feature=moderate type=rule control=sequence
//ff:what TestIsTrustedUser_Branches — covers empty specs and non-Author author branches of IsTrustedUser
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsTrustedUser_Branches(t *testing.T) {
	t.Run("EmptySpecs", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("author", &Author{TrustScore: 0.95})
		got, _ := IsTrustedUser(ctx, toulmin.Specs{})
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})

	t.Run("NonAuthor", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("author", "not-an-author")
		got, _ := IsTrustedUser(ctx, toulmin.Specs{&TrustScoreSpec{MinScore: 0.9}})
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})
}
