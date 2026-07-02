//ff:func feature=tangl type=engine control=sequence
//ff:what TestInitCompensation — tests that InitCompensation arms an empty compensation stack
package tangl

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestInitCompensation(t *testing.T) {
	ctx := toulmin.NewContext()
	InitCompensation(ctx)

	st := compensationStackOf(ctx)
	if st == nil {
		t.Fatal("expected compensation stack to be initialized, got nil")
	}
}
