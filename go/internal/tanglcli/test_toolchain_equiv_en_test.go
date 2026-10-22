//ff:func feature=cli type=command control=iteration dimension=1
//ff:what TestToolchainEquivEN — effects/gen output via the .geul path equals direct text-parse output
package tanglcli

import "testing"

// TestToolchainEquivEN pins the Phase 004 L4 equivalence (jointly with tangeul
// Phase 008): the effects and gen output produced from a .geul stream must
// match the output the old text-direct parse produced. Each fixture is checked by
// assertToolchainEquiv over both paths. americano is excluded — its undeclared
// `log` alias makes stream Save fail, so it has no .geul path.
func TestToolchainEquivEN(t *testing.T) {
	fixtures := []struct{ name, endpoint string }{
		{"transfer.md", "transfer"},
		{"access.md", "access"},
	}
	for _, f := range fixtures {
		assertToolchainEquiv(t, f.name, f.endpoint)
	}
}
