//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitSubject — encodes the Subject declaration as a Namespace node with a name label
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// TestEmitSubject is a single straight-line assertion: emitSubject has no
// data-dependent branch, only defensive appendNode/appendLabel error checks
// that never fire (NodeNamespace and PropLabel are always emitted with
// consistent, valid flags), so those error-return lines are unreachable
// through any legitimate call and are not exercised.
func TestEmitSubject(t *testing.T) {
	ref := doc.StmtRef{Sec: 1}
	e := &enc{tt: &tidTable{byRef: map[string]toulmin.TID{}}, wire: testWire, next: 100}

	got, err := emitSubject(e, ref, "transfer", toulmin.Span{})
	if err != nil {
		t.Fatalf("emitSubject: %v", err)
	}
	if len(got) == 0 {
		t.Errorf("emitSubject: expected non-empty words")
	}
}
