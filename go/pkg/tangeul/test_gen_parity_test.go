//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestGoldenGenParity — gen.Generate output is identical and gofmt-clean on both ASTs
package tangeul

import "testing"

// TestGoldenGenParity confirms the generated Go source is identical whether gen
// consumes the direct text-parse AST or the reloaded golden .geul AST, for each
// step-1 fixture. The per-fixture work is in checkGenParity.
func TestGoldenGenParity(t *testing.T) {
	for _, name := range []string{"transfer", "access"} {
		checkGenParity(t, name)
	}
}
