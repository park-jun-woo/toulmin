//ff:func feature=approve type=rule control=sequence
//ff:what TestHasApprovalRoleEdgeCases — HasApprovalRole covers the empty-specs and non-string approverRole branches
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// TestHasApprovalRoleEdgeCases covers the two branches of HasApprovalRole
// not exercised by TestHasApprovalRole's match/mismatch cases: an empty
// specs slice short-circuits to false, and a non-string (or missing)
// "approverRole" context value fails the type assertion and returns false.
func TestHasApprovalRoleEdgeCases(t *testing.T) {
	t.Run("empty specs returns false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("approverRole", "finance")
		got, val := HasApprovalRole(ctx, toulmin.Specs{})
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})

	t.Run("non-string approverRole returns false", func(t *testing.T) {
		ab := &ApproverSpec{Role: "finance"}
		ctx := toulmin.NewContext()
		ctx.Set("approverRole", 42)
		got, val := HasApprovalRole(ctx, toulmin.Specs{ab})
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})

	t.Run("missing approverRole returns false", func(t *testing.T) {
		ab := &ApproverSpec{Role: "finance"}
		ctx := toulmin.NewContext()
		got, val := HasApprovalRole(ctx, toulmin.Specs{ab})
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})
}
