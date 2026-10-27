//ff:func feature=tangl type=util control=sequence
//ff:what TestMetaRecord — wraps one Value as a Meta payload; an invalid Value fails
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
)

func TestMetaRecord(t *testing.T) {
	ws, err := metaRecord(core.MetaProfile, core.EncodeUintValue(7))
	if err != nil || len(ws) != 3 || ws[1] != 0x0882 {
		t.Fatalf("got %04X %v", ws, err)
	}
	if _, err := metaRecord(core.MetaProfile, core.Record{Kind: core.KindValue, Codec: 0x80}); err == nil {
		t.Fatal("invalid Value: want error")
	}
}
