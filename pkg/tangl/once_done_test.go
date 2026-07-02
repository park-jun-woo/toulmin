//ff:func feature=tangl type=engine control=sequence
//ff:what TestOnceDone — tests OnceDone across unset, wrong-type, true, and false stored values
package tangl

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestOnceDone(t *testing.T) {
	t.Run("UnsetKey", func(t *testing.T) {
		ctx := toulmin.NewContext()

		if OnceDone(ctx, "once:subj.case.node#0") {
			t.Fatal("expected OnceDone to be false for an unset key")
		}
	})

	t.Run("WrongType", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("once:subj.case.node#0", "not-a-bool")

		if OnceDone(ctx, "once:subj.case.node#0") {
			t.Fatal("expected OnceDone to be false for a non-bool value")
		}
	})

	t.Run("True", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("once:subj.case.node#0", true)

		if !OnceDone(ctx, "once:subj.case.node#0") {
			t.Fatal("expected OnceDone to be true when the stored bool is true")
		}
	})

	t.Run("False", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("once:subj.case.node#0", false)

		if OnceDone(ctx, "once:subj.case.node#0") {
			t.Fatal("expected OnceDone to be false when the stored bool is false")
		}
	})
}
