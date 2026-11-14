//ff:func feature=tangl type=util control=sequence
//ff:what TestRenderVerbatim — RenderVerbatim returns Doc.Source unchanged (byte identity) and nil for nil Doc
package doc

import (
	"bytes"
	"testing"
)

// TestRenderVerbatim verifies mode 1 decoding returns the Source bytes verbatim
// and that a nil Doc yields nil.
func TestRenderVerbatim(t *testing.T) {
	src := []byte("## tangl:Subject\n- this document is `t`\n")
	got := RenderVerbatim(&Doc{Source: src})
	if !bytes.Equal(got, src) {
		t.Errorf("RenderVerbatim = %q, want %q", got, src)
	}
	if RenderVerbatim(nil) != nil {
		t.Errorf("RenderVerbatim(nil) = non-nil, want nil")
	}
}
