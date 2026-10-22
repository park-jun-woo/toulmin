//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestGoldenRegen — regenerate/verify each step-1 golden .geul from its .md
package tangeul

import (
	"os"
	"testing"
)

// TestGoldenRegen is the golden regeneration harness for Phase 008 step 1. Run
// with TANGEUL_UPDATE=1 to rewrite pkg/tangeul/testdata/{transfer,access}.geul
// from the authored .md via parser.Parse -> assembleDoc -> Save; without the env
// it only verifies each golden reloads and matches the text-parse AST. The
// per-fixture work is in regenGolden. The americano exclusion is pinned separately
// by TestAmericanoExcluded.
func TestGoldenRegen(t *testing.T) {
	update := os.Getenv("TANGEUL_UPDATE") == "1"
	for _, name := range []string{"transfer", "access"} {
		regenGolden(t, name, update)
	}
}
