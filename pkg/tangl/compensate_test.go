//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompensate — tests Compensate's LIFO order, first-error stop, and no-init/empty-stack no-op paths
package tangl

import (
	"errors"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestCompensate(t *testing.T) {
	t.Run("LIFOOrder", func(t *testing.T) {
		ctx := toulmin.NewContext()
		InitCompensation(ctx)

		var order []string
		PushCompensation(ctx, func() error { order = append(order, "a"); return nil })
		PushCompensation(ctx, func() error { order = append(order, "b"); return nil })
		PushCompensation(ctx, func() error { order = append(order, "c"); return nil })

		if err := Compensate(ctx); err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		want := []string{"c", "b", "a"}
		if len(order) != len(want) {
			t.Fatalf("expected %v, got %v", want, order)
		}
		for i := range want {
			if order[i] != want[i] {
				t.Fatalf("expected %v, got %v", want, order)
			}
		}
	})

	t.Run("StopsAtFirstError", func(t *testing.T) {
		ctx := toulmin.NewContext()
		InitCompensation(ctx)

		boom := errors.New("boom")
		var ran []string
		// Armed in order: first (bottom) then second (top). Compensate runs LIFO,
		// so second's error must stop before first ever runs.
		PushCompensation(ctx, func() error { ran = append(ran, "first"); return nil })
		PushCompensation(ctx, func() error { ran = append(ran, "second"); return boom })

		err := Compensate(ctx)
		if !errors.Is(err, boom) {
			t.Fatalf("expected boom error, got %v", err)
		}
		if len(ran) != 1 || ran[0] != "second" {
			t.Fatalf("expected only 'second' to run before stopping, got %v", ran)
		}
	})

	t.Run("WithoutInitIsNoop", func(t *testing.T) {
		ctx := toulmin.NewContext()

		if err := Compensate(ctx); err != nil {
			t.Fatalf("expected nil error on uninitialized ctx, got %v", err)
		}
	})

	t.Run("EmptyStack", func(t *testing.T) {
		ctx := toulmin.NewContext()
		InitCompensation(ctx)

		if err := Compensate(ctx); err != nil {
			t.Fatalf("expected nil error for empty stack, got %v", err)
		}
	})
}
