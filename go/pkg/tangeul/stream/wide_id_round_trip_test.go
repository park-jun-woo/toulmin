//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestWideIDRoundTrip — a document whose Edge Ids need 2 words declares ID-WIDTH and reads back identically
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"reflect"

	"github.com/geul-org/geul/go/pkg/core"
)

// TestWideIDRoundTrip shifts every pass-1 Edge Id of the transfer document past
// 0xFFFF, encodes the BODY at a 2-word Edge Id width behind a HEAD that declares
// ID-WIDTH, and checks ReadDoc rebuilds the same AST and blocks — the declared
// width path end to end, without needing 65536 real records.
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
	body, last, err := emitBody(d, tt, toulmin.Wire{ID: 2, Span: 1})
	if err != nil || last <= 0xFFFF {
		t.Fatalf("emitBody: last=%d err=%v", last, err)
	}
	head, err := emitHead(d.Locale, 2)
	if err != nil {
		t.Fatal(err)
	}
	ws := core.WrapDocument(append(append(head, body...), 0x180F))
	got, err := ReadDoc(core.PutWords(ws))
	if err != nil {
		t.Fatalf("ReadDoc: %v", err)
	}
	if !reflect.DeepEqual(got.AST, d.AST) || !reflect.DeepEqual(got.Blocks, d.Blocks) {
		t.Fatalf("wide-id round trip drifted")
	}
}
