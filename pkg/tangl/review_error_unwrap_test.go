//ff:func feature=tangl type=model control=sequence
//ff:what TestReviewErrorUnwrap — tests ReviewError.Unwrap returns both wrapped errors
package tangl

import (
	"errors"
	"testing"
)

func TestReviewErrorUnwrap(t *testing.T) {
	cause := errors.New("cause-x")
	comp := errors.New("comp-y")
	r := &ReviewError{Cause: cause, Compensate: comp}

	errs := r.Unwrap()
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
	if errs[0] != cause {
		t.Fatalf("expected first error to be cause, got %v", errs[0])
	}
	if errs[1] != comp {
		t.Fatalf("expected second error to be compensation, got %v", errs[1])
	}
	if !errors.Is(r, cause) {
		t.Fatal("expected errors.Is to find cause via Unwrap")
	}
	if !errors.Is(r, comp) {
		t.Fatal("expected errors.Is to find compensation via Unwrap")
	}
}
