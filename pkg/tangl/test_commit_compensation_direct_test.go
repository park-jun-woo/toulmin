//ff:func feature=tangl type=engine control=sequence
//ff:what TestCommitCompensation_Direct — tests CommitCompensation sets the compensation key to nil
package tangl

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestCommitCompensation_Direct(t *testing.T) {
	ctx := toulmin.NewContext()
	CommitCompensation(ctx)

	v, ok := ctx.Get(compensationKey)
	if !ok {
		t.Fatal("expected compensationKey to be set")
	}
	if v != nil {
		t.Errorf("expected nil value, got %v", v)
	}
}
