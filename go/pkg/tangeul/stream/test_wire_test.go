//ff:type feature=tangl type=model
//ff:what testWire — the 1-word SPAN width stream unit tests encode and decode with
package stream

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// testWire is the narrowest document width, used by unit tests that build
// records by hand.
var testWire = toulmin.Wire{Span: 1}
