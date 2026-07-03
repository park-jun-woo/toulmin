//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestGoldenASTParity — golden .tangeul reload equals the direct text-parse AST up to Line
package tangeul

import "testing"

// TestGoldenASTParity asserts the core SSOT contract for step-1 fixtures: the AST
// recovered by reloading the golden .tangeul equals the AST the legacy text parser
// produces from the same .md, ignoring Line (unstored by the stream) and the
// top-level Document.Path (stream-nonstored metadata that Load stamps with the
// .tangeul path). Path is aligned here so EqualIgnoringLine can stay a strict
// Line-only comparison shared verbatim with tangl-en/tangl-ko.
func TestGoldenASTParity(t *testing.T) {
	for _, name := range []string{"transfer", "access"} {
		old, loaded := astPair(t, name)
		loaded.Path = old.Path
		if !EqualIgnoringLine(old, loaded) {
			t.Errorf("%s: reloaded AST differs from text-parse AST (ignoring Line)", name)
		}
	}
}
