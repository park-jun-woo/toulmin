//ff:func feature=state type=rule control=sequence
//ff:what TestIsExpired_EmptySpecs — tests IsExpired returns false when no specs provided
package state

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsExpired_EmptySpecs(t *testing.T) {
	ctx := toulmin.NewContext()
	got, evidence := IsExpired(ctx, toulmin.Specs{})
	if got != false {
		t.Errorf("expected false when specs is empty, got %v", got)
	}
	if evidence != nil {
		t.Errorf("expected nil evidence, got %v", evidence)
	}
}
