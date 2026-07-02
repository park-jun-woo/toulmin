//ff:func feature=approve type=rule control=sequence
//ff:what TestIsDirectManagerEdgeCases — IsDirectManager covers the missing-orgTree, missing-approverID, and missing-requesterID branches
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// TestIsDirectManagerEdgeCases covers the three type-assertion branches of
// IsDirectManager not exercised by TestIsDirectManager's is-manager/not-manager
// cases: a missing or wrong-type "orgTree", a missing or non-string
// "approverID", and a missing or non-string "requesterID" each short-circuit
// to false before OrgTree.IsDirectManager is called.
func TestIsDirectManagerEdgeCases(t *testing.T) {
	t.Run("wrong-type orgTree returns false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("orgTree", "not-an-org-tree")
		ctx.Set("approverID", "mgr-1")
		ctx.Set("requesterID", "emp-1")
		got, val := IsDirectManager(ctx, nil)
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})

	t.Run("missing orgTree returns false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("approverID", "mgr-1")
		ctx.Set("requesterID", "emp-1")
		got, _ := IsDirectManager(ctx, nil)
		if got {
			t.Errorf("got %v, want false", got)
		}
	})

	t.Run("non-string approverID returns false", func(t *testing.T) {
		org := &mockOrgTree{managers: map[string]string{"emp-1": "mgr-1"}}
		ctx := toulmin.NewContext()
		ctx.Set("orgTree", org)
		ctx.Set("approverID", 42)
		ctx.Set("requesterID", "emp-1")
		got, _ := IsDirectManager(ctx, nil)
		if got {
			t.Errorf("got %v, want false", got)
		}
	})

	t.Run("non-string requesterID returns false", func(t *testing.T) {
		org := &mockOrgTree{managers: map[string]string{"emp-1": "mgr-1"}}
		ctx := toulmin.NewContext()
		ctx.Set("orgTree", org)
		ctx.Set("approverID", "mgr-1")
		ctx.Set("requesterID", 42)
		got, _ := IsDirectManager(ctx, nil)
		if got {
			t.Errorf("got %v, want false", got)
		}
	})
}
