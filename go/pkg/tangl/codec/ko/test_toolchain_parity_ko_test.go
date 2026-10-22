//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestToolchainParityKo — validate/Lint/effects/gen agree on the KO- and EN-encoded forms of each fixture
package ko

import "testing"

// TestToolchainParityKo is the core contract of plans/tangl-ko §7: a KO-authored
// surface, once absorbed to a Doc, must drive the locale-blind toolchain (validate
// / Lint / effects / gen) to exactly the same result as its EN counterpart. Each
// shared fixture is compared by compareToolchainKo; because gen and effects are
// Line-independent, the only differences the shared AST could leak are none — that
// is the property under test.
func TestToolchainParityKo(t *testing.T) {
	for _, name := range []string{"americano", "access", "transfer", "rules"} {
		compareToolchainKo(t, name)
	}
}
