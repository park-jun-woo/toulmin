//ff:func feature=state type=rule control=sequence
//ff:what TestIsOwner_MissingKey — tests IsOwner when context keys are absent
package state

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsOwner_MissingKey(t *testing.T) {
	ctx := toulmin.NewContext()
	got, evidence := IsOwner(ctx, nil)
	if got != true {
		t.Errorf("expected true when both keys are absent (nil == nil), got %v", got)
	}
	if evidence != nil {
		t.Errorf("expected nil evidence, got %v", evidence)
	}
}
