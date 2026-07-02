//ff:func feature=tangl type=engine control=sequence
//ff:what TestReview — tests that Review wraps both errors into a ReviewError and logs them to stderr
package tangl

import (
	"bytes"
	"errors"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestReview(t *testing.T) {
	t.Run("WrapsBothErrors", func(t *testing.T) {
		ctx := toulmin.NewContext()
		cause := errors.New("withdraw failed downstream")
		comp := errors.New("refund failed")

		err := captureStderr(t, func() error {
			return Review(ctx, cause, comp)
		})

		if !errors.Is(err, cause) {
			t.Fatalf("expected errors.Is(err, cause) to hold, got %v", err)
		}
		if !errors.Is(err, comp) {
			t.Fatalf("expected errors.Is(err, comp) to hold, got %v", err)
		}

		var reviewErr *ReviewError
		if !errors.As(err, &reviewErr) {
			t.Fatalf("expected errors.As to find *ReviewError, got %v", err)
		}
		if reviewErr.Cause != cause || reviewErr.Compensate != comp {
			t.Fatalf("expected ReviewError to preserve both errors, got %+v", reviewErr)
		}
	})

	t.Run("LogsToStderr", func(t *testing.T) {
		ctx := toulmin.NewContext()
		cause := errors.New("cause-x")
		comp := errors.New("comp-y")

		logged := captureStderrOutput(t, func() error {
			return Review(ctx, cause, comp)
		})

		if !bytes.Contains(logged, []byte("cause-x")) || !bytes.Contains(logged, []byte("comp-y")) {
			t.Fatalf("expected stderr log to mention both errors, got %q", logged)
		}
	})
}
