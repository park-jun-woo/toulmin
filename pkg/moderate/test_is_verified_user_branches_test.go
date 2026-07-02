//ff:func feature=moderate type=rule control=sequence
//ff:what TestIsVerifiedUser_NonAuthor — covers non-Author author branch of IsVerifiedUser
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsVerifiedUser_NonAuthor(t *testing.T) {
	ctx := toulmin.NewContext()
	ctx.Set("author", "not-an-author")
	got, _ := IsVerifiedUser(ctx, nil)
	if got != false {
		t.Errorf("got %v, want false", got)
	}
}
