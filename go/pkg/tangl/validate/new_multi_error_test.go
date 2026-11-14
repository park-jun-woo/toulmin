//ff:func feature=tangl type=util control=sequence
//ff:what TestNewMultiError — tests newMultiError for empty-slice-nil and non-empty-wrap branches
package validate

import (
	"errors"
	"testing"
)

func TestNewMultiError(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		if err := newMultiError(nil); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
		if err := newMultiError([]error{}); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("NonEmpty", func(t *testing.T) {
		errs := []error{errors.New("boom")}
		err := newMultiError(errs)
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		me, ok := err.(*multiError)
		if !ok {
			t.Fatalf("expected *multiError, got %T", err)
		}
		if len(me.errs) != 1 || me.errs[0].Error() != "boom" {
			t.Fatalf("expected wrapped errs to match input, got %v", me.errs)
		}
	})
}
