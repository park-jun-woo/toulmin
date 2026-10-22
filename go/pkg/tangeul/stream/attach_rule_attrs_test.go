//ff:func feature=tangl type=util control=sequence
//ff:what TestAttachRuleAttrs — host function, with and grounds set their flags and lists; each issue failure propagates
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestAttachRuleAttrs(t *testing.T) {
	var pre, body []core.Word
	var p toulmin.NodePacket
	if err := attachRuleAttrs(testEnc(nil), &pre, &body, &ast.Ref{Name: "f"}, []string{"w"}, []string{"g1", "g2"}, toulmin.Span{}, &p); err != nil {
		t.Fatal(err)
	}
	want := toulmin.NodeFlagFunc | toulmin.NodeFlagWith | toulmin.NodeFlagGrounds
	if p.Flags != want || p.Func == nil || len(p.With) != 1 || len(p.Grounds) != 2 {
		t.Fatalf("got %+v", p)
	}
	var none toulmin.NodePacket
	if err := attachRuleAttrs(testEnc(nil), &pre, &body, nil, nil, nil, toulmin.Span{}, &none); err != nil || none.Flags != 0 {
		t.Fatalf("no attrs: %+v %v", none, err)
	}
	// Issues: FUNC name TEXT, FUNC node, LABEL, with TEXT, grounds TEXT.
	failsAtEachID(t, 5, nil, func(e *enc) error {
		var q toulmin.NodePacket
		return attachRuleAttrs(e, &pre, &body, &ast.Ref{Name: "f"}, []string{"w"}, []string{"g"}, toulmin.Span{}, &q)
	})
}
