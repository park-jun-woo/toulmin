//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitHead — emitHead's HEAD reads back through readHead with the same width and locale
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/core"
)

// TestEmitHead round-trips the HEAD for the default width without a locale and
// for a 4-word width with a locale, and rejects an undeclarable width.
func TestEmitHead(t *testing.T) {
	for _, c := range []struct {
		w      int
		locale string
	}{{1, ""}, {4, "ko"}} {
		ws, err := emitHead(c.locale, c.w)
		if err != nil {
			t.Fatal(err)
		}
		doc := core.WrapDocument(append(ws, 0x180F))
		gd, _ := core.SplitDocument(doc)
		h, err := readHead(gd.Head)
		if err != nil || h.idWidth != c.w || h.locale != c.locale {
			t.Fatalf("%+v: got %+v %v", c, h, err)
		}
	}
	if _, err := emitHead("", 3); err == nil {
		t.Fatal("width 3: want error")
	}
}
