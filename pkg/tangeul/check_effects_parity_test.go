//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what checkEffectsParity — assert effects.Closure matches across routes for one fixture's endpoints
package tangeul

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/effects"
)

// checkEffectsParity computes effects.Closure over both the text-parse AST and the
// reloaded golden .tangeul AST for every provides endpoint of the named fixture
// and fails on any difference. effects.Entry carries no Line or Path, so the
// comparison is exact (no normalization) — the retargeting guarantee for the
// do/undo execution plan.
func checkEffectsParity(t *testing.T, name string) {
	t.Helper()
	old, loaded := astPair(t, name)
	for _, ep := range old.Provides {
		wantCl, wantErr := effects.Closure(old, ep.Name)
		gotCl, gotErr := effects.Closure(loaded, ep.Name)
		if (wantErr == nil) != (gotErr == nil) {
			t.Errorf("%s/%s: Closure error mismatch: old=%v new=%v", name, ep.Name, wantErr, gotErr)
		}
		if !reflect.DeepEqual(wantCl, gotCl) {
			t.Errorf("%s/%s: closure differs\nold=%+v\nnew=%+v", name, ep.Name, wantCl, gotCl)
		}
	}
}
