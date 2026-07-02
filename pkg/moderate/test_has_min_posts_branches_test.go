//ff:func feature=moderate type=rule control=sequence
//ff:what TestHasMinPosts_Branches — covers empty specs and non-Author author branches of HasMinPosts
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestHasMinPosts_Branches(t *testing.T) {
	t.Run("EmptySpecs", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("author", &Author{PostCount: 100})
		got, _ := HasMinPosts(ctx, toulmin.Specs{})
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})

	t.Run("NonAuthor", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("author", "not-an-author")
		got, _ := HasMinPosts(ctx, toulmin.Specs{&MinPostsSpec{MinPosts: 10}})
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})
}
