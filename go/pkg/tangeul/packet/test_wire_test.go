//ff:type feature=tangl type=model
//ff:what testWire — the 1-word Edge Id / 1-word SPAN widths the packet tests encode with
package packet

import "github.com/geul-org/geul/go/pkg/toulmin"

// testWire is the narrowest document widths (1-word Edge Id, 1-word SPAN
// fields), the common case for small TANGEUL documents.
var testWire = toulmin.Wire{ID: 1, Span: 1}
