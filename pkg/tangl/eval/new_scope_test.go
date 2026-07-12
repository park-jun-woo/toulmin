//ff:func feature=tangl type=engine control=sequence
//ff:what TestNewScope — tests fresh-scope seeding of ctx, specs, and empty binders
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestNewScope(t *testing.T) {
	ctx := testCtx()
	specs := toulmin.Specs{&testSpec{Name: "s"}}
	sc := newScope(ctx, specs)
	if sc.ctx != toulmin.Context(ctx) {
		t.Fatal("scope ctx not seeded")
	}
	if len(sc.specs) != 1 {
		t.Fatalf("scope specs not seeded: %v", sc.specs)
	}
	if sc.binders == nil || len(sc.binders) != 0 {
		t.Fatalf("scope binders should start empty, got %v", sc.binders)
	}
}
