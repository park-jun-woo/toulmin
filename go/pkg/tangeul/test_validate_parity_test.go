//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestGoldenValidateParity — validate/Lint agree between routes after path:line normalization
package tangeul

import "testing"

// TestGoldenValidateParity confirms validate.Validate and validate.Lint report the
// same violations for the direct text-parse AST and the reloaded golden .geul
// AST, for each step-1 fixture. The per-fixture work is in checkValidateParity.
func TestGoldenValidateParity(t *testing.T) {
	for _, name := range []string{"transfer", "access"} {
		checkValidateParity(t, name)
	}
}
