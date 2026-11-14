//ff:func feature=tangl type=codegen control=sequence
//ff:what TestModuleLine — covers a present module directive and its absence
package gen

import "testing"

func TestModuleLine(t *testing.T) {
	t.Run("module directive found", func(t *testing.T) {
		got, ok := moduleLine([]byte("// header\n\n  module   example.com/x  \n\ngo 1.18\n"))
		if !ok || got != "example.com/x" {
			t.Fatalf("got %q ok=%v, want example.com/x true", got, ok)
		}
	})

	t.Run("no module directive", func(t *testing.T) {
		if got, ok := moduleLine([]byte("go 1.18\nrequire x v1.0.0\n")); ok {
			t.Fatalf("expected ok=false, got %q true", got)
		}
	})
}
