//ff:type feature=tangl type=model
//ff:what testWire — the 1-word SPAN width the packet tests encode with
package packet

import "github.com/geul-org/gogeul/pkg/toulmin"

// testWire is the narrowest document width (1-word SPAN fields), the common
// case for small TANGEUL documents.
var testWire = toulmin.Wire{Span: 1}
