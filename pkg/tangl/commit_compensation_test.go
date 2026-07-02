//ff:func feature=tangl type=engine control=sequence
//ff:what TestCommitCompensation — tests that CommitCompensation makes Compensate a no-op and tolerates a missing InitCompensation call
package tangl

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestCommitCompensation(t *testing.T) {
	t.Run("MakesCompensateNoop", func(t *testing.T) {
		ctx := toulmin.NewContext()
		InitCompensation(ctx)

		ran := false
		PushCompensation(ctx, func() error { ran = true; return nil })

		CommitCompensation(ctx)

		if err := Compensate(ctx); err != nil {
			t.Fatalf("expected nil error after commit, got %v", err)
		}
		if ran {
			t.Fatal("expected committed compensation to not run")
		}
	})

	t.Run("WithoutInitDoesNotPanic", func(t *testing.T) {
		ctx := toulmin.NewContext()

		CommitCompensation(ctx)
	})
}
