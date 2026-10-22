//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestClassifyRecord — each record family gets its reader class; foreign records fail
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
)

func TestClassifyRecord(t *testing.T) {
	iss := func(f uint16) core.Record {
		return core.Record{Kind: core.KindEdge, Codec: core.CodecIssuer1Word, Issuer: []core.Word{0}, Format: f}
	}
	cases := []struct {
		r    core.Record
		want recClass
	}{
		{core.Record{Kind: core.KindEdge, Codec: toulmin.CodecCaseDeclare}, classNode},
		{core.Record{Kind: core.KindEdge, Codec: toulmin.CodecDo}, classLink},
		{core.Record{Kind: core.KindEdge, Codec: toulmin.CodecPredicateDeclare}, classPayload},
		{iss(toulmin.FormatText), classPayload},
		{iss(toulmin.FormatProse), classControl},
		{iss(toulmin.FormatSource), classSource},
	}
	for _, c := range cases {
		if got, err := classifyRecord(c.r); err != nil || got != c.want {
			t.Fatalf("%+v: %d %v", c.r, got, err)
		}
	}
	for _, bad := range []core.Record{{Kind: core.KindMeta}, iss(0xCFFF)} {
		if _, err := classifyRecord(bad); err == nil {
			t.Fatalf("%+v: want error", bad)
		}
	}
}
