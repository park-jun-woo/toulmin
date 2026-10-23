//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestWideIDRoundTrip — Edge Ids past 0xFFFF encode as 2-word LOCAL References and read back identically
package stream

import (
	"reflect"
	"testing"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// TestWideIDRoundTrip shifts every pass-1 Edge Id of the transfer document past
// 0xFFFF, encodes the BODY (each LOCAL Reference then takes 2 id words) and
// checks ReadDoc rebuilds the same AST and blocks — the wide-id path end to
// end, without needing 65536 real records.
func TestWideIDRoundTrip(t *testing.T) {
	d := buildTransferDoc()
	tt, err := assignTIDs(d)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range []map[string]toulmin.TID{tt.byRef, tt.aliasByName, tt.caseByName, tt.nodeByCase} {
		for k := range m {
			m[k] += 0x10000
		}
	}
	tt.primaryEnd += 0x10000
	body, err := emitBody(d, tt, testWire)
	if err != nil {
		t.Fatalf("emitBody: %v", err)
	}
	ws := core.WrapDocument(append(append(emitHead(d.Locale), body...), 0x180F))
	got, err := ReadDoc(core.PutWords(ws))
	if err != nil {
		t.Fatalf("ReadDoc: %v", err)
	}
	if !reflect.DeepEqual(got.AST, d.AST) || !reflect.DeepEqual(got.Blocks, d.Blocks) {
		t.Fatalf("wide-id round trip drifted")
	}
}
