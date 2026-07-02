//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompensationStackOf — tests compensationStackOf for missing key, wrong type, and valid stack
package tangl

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestCompensationStackOf(t *testing.T) {
	// key not set: !ok branch -> nil
	ctx1 := toulmin.NewContext()
	if st := compensationStackOf(ctx1); st != nil {
		t.Errorf("expected nil for missing key, got %+v", st)
	}

	// key set to wrong type: type assertion failure branch -> nil
	ctx2 := toulmin.NewContext()
	ctx2.Set(compensationKey, "not-a-stack")
	if st := compensationStackOf(ctx2); st != nil {
		t.Errorf("expected nil for wrong type, got %+v", st)
	}

	// key set to correct type: returns the stack
	ctx3 := toulmin.NewContext()
	want := &compensationStack{}
	ctx3.Set(compensationKey, want)
	if st := compensationStackOf(ctx3); st != want {
		t.Errorf("expected %p, got %p", want, st)
	}
}
