//ff:func feature=tangl type=util control=sequence
//ff:what TestReadHead — readHead requires both stamps and propagates a bad record
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/core"
)

// TestReadHead checks an empty HEAD (no stamps) and a HEAD with an unsupported
// Meta are rejected.
func TestReadHead(t *testing.T) {
	if _, err := readHead(nil); err == nil {
		t.Fatal("no stamps: want error")
	}
	if _, err := readHead([]core.Record{{Kind: core.KindMeta, Codec: 0x12}}); err == nil {
		t.Fatal("unsupported Meta: want error")
	}
}
