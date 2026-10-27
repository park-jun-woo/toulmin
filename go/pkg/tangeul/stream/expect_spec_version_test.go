//ff:func feature=tangl type=util control=sequence
//ff:what TestExpectSpecVersion — this build's version passes; another version and a malformed payload fail
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/toulmin"
)

func TestExpectSpecVersion(t *testing.T) {
	if err := expectSpecVersion(core.EncodeSpecVersion(toulmin.SpecVersion)); err != nil {
		t.Fatal(err)
	}
	if err := expectSpecVersion(core.EncodeSpecVersion(2)); err == nil {
		t.Fatal("other version: want error")
	}
	if err := expectSpecVersion(core.Record{Kind: core.KindMeta, Codec: core.MetaSpecVersion}); err == nil {
		t.Fatal("empty: want error")
	}
}
