//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitHead — emitHead's HEAD reads back through readHead with the same locale and declares no width
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
)

// TestEmitHead round-trips the HEAD without and with a locale and checks it
// holds only SPEC-VERSION, PROFILE and the optional DEFAULT-LANGUAGE.
func TestEmitHead(t *testing.T) {
	for _, c := range []struct {
		locale string
		metas  int
	}{{"", 2}, {"ko", 3}} {
		doc := core.WrapDocument(append(emitHead(c.locale), 0x180F))
		gd, _ := core.SplitDocument(doc)
		h, err := readHead(gd.Head)
		if err != nil || h.locale != c.locale || len(gd.Head) != c.metas {
			t.Fatalf("%+v: got %+v %d %v", c, h, len(gd.Head), err)
		}
	}
}
