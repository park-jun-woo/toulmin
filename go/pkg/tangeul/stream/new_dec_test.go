//ff:func feature=tangl type=util control=sequence
//ff:what TestNewDec — newDec initializes empty maps, textNext=1, an empty AST and curCase=-1
package stream

import "testing"

// TestNewDec checks the fresh Reader state newDec builds: every resolution
// map is non-nil and empty, the AST is a fresh empty Document, textNext
// starts at 1 and curCase starts outside any case (-1).
func TestNewDec(t *testing.T) {
	d := newDec(testWire)
	if d.ast == nil {
		t.Fatalf("ast = nil")
	}
	if d.ast.Subject != "" || len(d.ast.Cases) != 0 {
		t.Errorf("ast = %+v, want a fresh empty Document", d.ast)
	}
	if d.texts == nil || len(d.texts) != 0 {
		t.Errorf("texts = %v, want empty non-nil map", d.texts)
	}
	if d.wire != testWire || d.seen == nil || len(d.seen) != 0 {
		t.Errorf("wire/seen = %+v/%v, want testWire and an empty non-nil map", d.wire, d.seen)
	}
	if d.nodeName == nil || len(d.nodeName) != 0 {
		t.Errorf("nodeName = %v, want empty non-nil map", d.nodeName)
	}
	if d.nodeType == nil || len(d.nodeType) != 0 {
		t.Errorf("nodeType = %v, want empty non-nil map", d.nodeType)
	}
	if d.funcRef == nil || len(d.funcRef) != 0 {
		t.Errorf("funcRef = %v, want empty non-nil map", d.funcRef)
	}
	if d.caseByTID == nil || len(d.caseByTID) != 0 {
		t.Errorf("caseByTID = %v, want empty non-nil map", d.caseByTID)
	}
	if d.epByTID == nil || len(d.epByTID) != 0 {
		t.Errorf("epByTID = %v, want empty non-nil map", d.epByTID)
	}
	if d.trByTID == nil || len(d.trByTID) != 0 {
		t.Errorf("trByTID = %v, want empty non-nil map", d.trByTID)
	}
	if d.defByTID == nil || len(d.defByTID) != 0 {
		t.Errorf("defByTID = %v, want empty non-nil map", d.defByTID)
	}
	if d.curCase != -1 {
		t.Errorf("curCase = %d, want -1", d.curCase)
	}
}
