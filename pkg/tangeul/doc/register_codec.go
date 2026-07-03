//ff:func feature=tangl type=util control=sequence
//ff:what RegisterCodec — register a Codec by locale; duplicate locale registration panics
package doc

import "fmt"

// codecRegistry maps a locale tag to its registered Codec. The CLI assembly point
// (tanglcli) registers the EN/KO codecs explicitly; there is no init().
var codecRegistry = map[string]Codec{}

// RegisterCodec registers c under its Locale(). Registering a locale twice is an
// assembly bug and panics (README contract §3).
func RegisterCodec(c Codec) {
	locale := c.Locale()
	if _, ok := codecRegistry[locale]; ok {
		panic(fmt.Sprintf("tangeul/doc: duplicate codec registration for locale %q", locale))
	}
	codecRegistry[locale] = c
}
