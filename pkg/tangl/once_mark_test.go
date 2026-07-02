//ff:func feature=tangl type=engine control=sequence
//ff:what TestOnceMark — tests that OnceMark flips OnceDone for its own key while leaving other keys independent
package tangl

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestOnceMark(t *testing.T) {
	t.Run("ThenDone", func(t *testing.T) {
		ctx := toulmin.NewContext()
		key := "once:subj.case.node#0"

		if OnceDone(ctx, key) {
			t.Fatal("expected OnceDone to be false before OnceMark")
		}

		OnceMark(ctx, key)

		if !OnceDone(ctx, key) {
			t.Fatal("expected OnceDone to be true after OnceMark")
		}
	})

	t.Run("KeysAreIndependent", func(t *testing.T) {
		ctx := toulmin.NewContext()

		OnceMark(ctx, "once:subj.case.a#0")

		if OnceDone(ctx, "once:subj.case.b#0") {
			t.Fatal("expected an unrelated once key to remain unmarked")
		}
	})
}
