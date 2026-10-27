//ff:func feature=tangl type=util control=sequence
//ff:what TestToWireSpan — toWireSpan widens offset and length without a width limit
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// TestToWireSpan checks offsets past the old 16-bit limit pass through: the
// width is declared per document, so only the record encoder enforces it.
func TestToWireSpan(t *testing.T) {
	if got := toWireSpan(doc.Span{Off: 0x12345, Len: 0x10000}); got != (toulmin.Span{Off: 0x12345, Len: 0x10000}) {
		t.Fatalf("got %+v", got)
	}
}
