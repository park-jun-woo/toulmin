//ff:func feature=tangl type=util control=sequence
//ff:what TestMultiErrorError — tests multiError.Error for empty, single, and multiple error branches
package validate

import (
	"errors"
	"testing"
)

func TestMultiErrorError(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		m := &multiError{}
		if got := m.Error(); got != "" {
			t.Fatalf("expected empty string, got %q", got)
		}
	})

	t.Run("Single", func(t *testing.T) {
		m := &multiError{errs: []error{errors.New("one")}}
		if got := m.Error(); got != "one" {
			t.Fatalf("expected %q, got %q", "one", got)
		}
	})

	t.Run("Multiple", func(t *testing.T) {
		m := &multiError{errs: []error{errors.New("one"), errors.New("two")}}
		want := "one\ntwo"
		if got := m.Error(); got != want {
			t.Fatalf("expected %q, got %q", want, got)
		}
	})
}
