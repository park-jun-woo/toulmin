//ff:func feature=tangl type=model control=sequence
//ff:what TestReviewErrorError — tests ReviewError.Error formatting with cause and compensation errors, including nil fields
package tangl

import (
	"errors"
	"strings"
	"testing"
)

func TestReviewErrorError(t *testing.T) {
	t.Run("CauseAndCompensate", func(t *testing.T) {
		cause := errors.New("cause-x")
		comp := errors.New("comp-y")
		r := &ReviewError{Cause: cause, Compensate: comp}

		got := r.Error()
		if !strings.Contains(got, "cause-x") {
			t.Fatalf("expected error message to contain cause, got %q", got)
		}
		if !strings.Contains(got, "comp-y") {
			t.Fatalf("expected error message to contain compensation error, got %q", got)
		}
		if !strings.Contains(got, "REVIEW required") {
			t.Fatalf("expected error message to contain REVIEW marker, got %q", got)
		}
	})

	t.Run("NilFields", func(t *testing.T) {
		r := &ReviewError{}

		got := r.Error()
		if !strings.Contains(got, "<nil>") {
			t.Fatalf("expected error message to contain <nil> for unset fields, got %q", got)
		}
	})
}
