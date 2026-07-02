//ff:func feature=policy type=rule control=sequence
//ff:what TestHasHeader_Branches — covers empty specs and headers type-assertion failure branches of HasHeader
package policy

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestHasHeader_Branches(t *testing.T) {
	t.Run("empty specs", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("headers", map[string]string{"X-Internal-Token": "secret"})

		got, evidence := HasHeader(ctx, toulmin.Specs{})
		if got {
			t.Errorf("expected false for empty specs, got %v", got)
		}
		if evidence != nil {
			t.Errorf("expected nil evidence, got %v", evidence)
		}
	})

	t.Run("headers wrong type", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("headers", "not-a-map")

		got, _ := HasHeader(ctx, toulmin.Specs{&HeaderSpec{Header: "X-Internal-Token"}})
		if got {
			t.Errorf("expected false when headers is not a map[string]string, got %v", got)
		}
	})

	t.Run("headers missing", func(t *testing.T) {
		ctx := toulmin.NewContext()

		got, _ := HasHeader(ctx, toulmin.Specs{&HeaderSpec{Header: "X-Internal-Token"}})
		if got {
			t.Errorf("expected false when headers is unset, got %v", got)
		}
	})
}
