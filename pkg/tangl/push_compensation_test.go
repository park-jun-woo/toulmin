//ff:func feature=tangl type=engine control=sequence
//ff:what TestPushCompensation — tests PushCompensation's lazy-init, nil-ignore, and append-to-existing-stack behavior
package tangl

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestPushCompensation(t *testing.T) {
	t.Run("WithoutInitDoesNotPanic", func(t *testing.T) {
		ctx := toulmin.NewContext()

		ran := false
		PushCompensation(ctx, func() error { ran = true; return nil })

		if err := Compensate(ctx); err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !ran {
			t.Fatal("expected lazily-armed compensation to run")
		}
	})

	t.Run("IgnoresNilFunc", func(t *testing.T) {
		ctx := toulmin.NewContext()
		InitCompensation(ctx)

		PushCompensation(ctx, nil)

		if err := Compensate(ctx); err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	})

	t.Run("AppendsToExistingStack", func(t *testing.T) {
		ctx := toulmin.NewContext()
		InitCompensation(ctx)

		firstRan := false
		secondRan := false
		PushCompensation(ctx, func() error { firstRan = true; return nil })
		PushCompensation(ctx, func() error { secondRan = true; return nil })

		if err := Compensate(ctx); err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !firstRan || !secondRan {
			t.Fatalf("expected both compensations to run, got first=%v second=%v", firstRan, secondRan)
		}
	})
}
