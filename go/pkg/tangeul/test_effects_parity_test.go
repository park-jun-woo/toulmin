//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestGoldenEffectsParity — effects.Closure agrees between text-parse and reloaded AST
package tangeul

import "testing"

// TestGoldenEffectsParity confirms the effect closure of every provides endpoint
// is identical whether computed over the direct text-parse AST or the reloaded
// golden .geul AST, for each step-1 fixture. The per-fixture work is in
// checkEffectsParity.
func TestGoldenEffectsParity(t *testing.T) {
	for _, name := range []string{"transfer", "access"} {
		checkEffectsParity(t, name)
	}
}
