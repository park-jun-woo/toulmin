//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEncoderGoldenXCheck — the tangl-en encoder cross-checks against every step-1 golden
package tangeul

import "testing"

// TestEncoderGoldenXCheck is the Phase 008 step-2 cross-verification: for each step-1
// fixture, the stream the tangl-en encoder now produces must be structurally equivalent
// to the hand-assembled golden frozen in step 1 — same AST up to Line and the same
// ordered sequence of statement/section blocks (StmtRef coordinates). This proves the
// two independent stream constructions agree, without conflating them: the golden stays
// a hand-built stand-in (placeholder prose spans) so the check is not tautological, and
// byte-level identity of the encoder's real prose is verified elsewhere (TestEncodeDecodeEN).
// americano is excluded — it has no golden (its undeclared `log` alias fails stream Save).
// Per-fixture work is in checkEncoderXCheck.
func TestEncoderGoldenXCheck(t *testing.T) {
	for _, name := range []string{"transfer", "access"} {
		checkEncoderXCheck(t, name)
	}
}
