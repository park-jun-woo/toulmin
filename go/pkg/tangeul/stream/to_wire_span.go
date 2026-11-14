//ff:func feature=tangl type=util control=sequence
//ff:what toWireSpan — convert a logical doc.Span to a wire toulmin.Span
package stream

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// toWireSpan converts a logical doc.Span (offset/length into Source) to a wire
// toulmin.Span. The wire width is declared per document by the SOURCE record, so
// no width limit applies here; the record encoder rejects a span the declared
// width cannot hold.
func toWireSpan(s doc.Span) toulmin.Span {
	return toulmin.Span{Off: uint64(s.Off), Len: uint64(s.Len)}
}
