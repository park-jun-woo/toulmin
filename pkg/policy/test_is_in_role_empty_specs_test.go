//ff:func feature=policy type=rule control=sequence
//ff:what TestIsInRole_EmptySpecs — covers IsInRole len(specs)==0 branch
package policy

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsInRole_EmptySpecs(t *testing.T) {
	ctx := toulmin.NewContext()
	ctx.Set("user", &testUser{Role: "admin"})
	ctx.Set("role", "admin")

	got, evidence := IsInRole(ctx, toulmin.Specs{})
	if got {
		t.Errorf("expected false for empty specs, got %v", got)
	}
	if evidence != nil {
		t.Errorf("expected nil evidence, got %v", evidence)
	}
}
