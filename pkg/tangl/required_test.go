//ff:func feature=tangl type=engine control=sequence
//ff:what TestRequired — tests Required across missing, nil, all-present, and no-field cases
package tangl

import (
	"errors"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestRequired(t *testing.T) {
	t.Run("MissingField", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("from account", "acc-1")

		err := Required(ctx, "from account", "amount")
		if err == nil {
			t.Fatal("expected error for missing field, got nil")
		}
		if !errors.Is(err, ErrRequired) {
			t.Fatalf("expected errors.Is(err, ErrRequired), got %v", err)
		}
	})

	t.Run("NilField", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("amount", nil)

		if err := Required(ctx, "amount"); !errors.Is(err, ErrRequired) {
			t.Fatalf("expected ErrRequired for nil field, got %v", err)
		}
	})

	t.Run("AllPresent", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("from account", "acc-1")
		ctx.Set("amount", 100)

		if err := Required(ctx, "from account", "amount"); err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	})

	t.Run("NoFields", func(t *testing.T) {
		ctx := toulmin.NewContext()

		if err := Required(ctx); err != nil {
			t.Fatalf("expected nil error for empty field list, got %v", err)
		}
	})
}
