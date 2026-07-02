//ff:func feature=policy type=rule control=sequence
//ff:what TestIsIPInList_Branches — covers empty specs and clientIP type-assertion failure branches of IsIPInList
package policy

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsIPInList_Branches(t *testing.T) {
	t.Run("empty specs", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("clientIP", "1.2.3.4")

		got, evidence := IsIPInList(ctx, toulmin.Specs{})
		if got {
			t.Errorf("expected false for empty specs, got %v", got)
		}
		if evidence != nil {
			t.Errorf("expected nil evidence, got %v", evidence)
		}
	})

	t.Run("clientIP wrong type", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("clientIP", 12345)

		list := &IPListSpec{Purpose: "blocklist", List: []string{"1.2.3.4"}}
		got, _ := IsIPInList(ctx, toulmin.Specs{list})
		if got {
			t.Errorf("expected false when clientIP is not a string, got %v", got)
		}
	})

	t.Run("clientIP unset", func(t *testing.T) {
		ctx := toulmin.NewContext()

		list := &IPListSpec{Purpose: "blocklist", List: []string{"1.2.3.4"}}
		got, _ := IsIPInList(ctx, toulmin.Specs{list})
		if got {
			t.Errorf("expected false when clientIP is unset, got %v", got)
		}
	})
}
