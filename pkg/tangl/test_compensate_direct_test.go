//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompensate_Direct — tests Compensate nil-stack, error-stop, and success paths directly
package tangl

import (
	"errors"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestCompensate_Direct(t *testing.T) {
	// nil stack (never initialized): st == nil branch
	ctx1 := toulmin.NewContext()
	if err := Compensate(ctx1); err != nil {
		t.Fatalf("expected nil error for uninitialized ctx, got %v", err)
	}

	// non-nil stack, function returns error: loop error-return branch
	ctx2 := toulmin.NewContext()
	InitCompensation(ctx2)
	boom := errors.New("boom")
	PushCompensation(ctx2, func() error { return boom })
	if err := Compensate(ctx2); !errors.Is(err, boom) {
		t.Fatalf("expected boom error, got %v", err)
	}

	// non-nil stack, all functions succeed: loop completes, final nil return
	ctx3 := toulmin.NewContext()
	InitCompensation(ctx3)
	ran := false
	PushCompensation(ctx3, func() error { ran = true; return nil })
	if err := Compensate(ctx3); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !ran {
		t.Fatal("expected compensation function to run")
	}
}
