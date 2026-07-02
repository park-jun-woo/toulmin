//ff:func feature=policy type=rule control=sequence dimension=1
//ff:what TestIsOwner_UserUnset — covers IsOwner when the "user" key is never set on the context
package policy

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsOwner_UserUnset(t *testing.T) {
	ctx := toulmin.NewContext()
	got, evidence := IsOwner(ctx, toulmin.Specs{&OwnerSpec{}})
	if got {
		t.Errorf("expected false when user key is unset, got %v", got)
	}
	if evidence != nil {
		t.Errorf("expected nil evidence, got %v", evidence)
	}
}
