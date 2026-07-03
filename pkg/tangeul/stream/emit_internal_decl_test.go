//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitInternalDecl — on/every heads, an until edge to a declared case, and an undeclared until error
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestEmitInternalDecl exercises emitInternalDecl's Kind dispatch (OnEvent vs
// EveryTick, which selects Event vs Interval as the head label source) and the
// `until` edge: absent, resolved against a declared case, and referencing an
// undeclared case (a hard error). appendNode/appendLabel/appendEdge never fail
// here (NodeTrigger/PropUntil are always emitted with consistent flags), so
// their error-return lines are unreachable through any legitimate call and are
// not exercised.
func TestEmitInternalDecl(t *testing.T) {
	ref := doc.StmtRef{Sec: 7, Path: []int{0}}
	span := word.Span{}

	cases := []struct {
		name    string
		in      ast.Internal
		wantErr bool
	}{
		{"on event", ast.Internal{Kind: ast.OnEvent, Event: "sensor.buttonPress"}, false},
		{"every tick, no until", ast.Internal{Kind: ast.EveryTick, Interval: "1s"}, false},
		{"every tick, declared until", ast.Internal{Kind: ast.EveryTick, Interval: "1s", Until: "target"}, false},
		{"every tick, undeclared until", ast.Internal{Kind: ast.EveryTick, Interval: "1s", Until: "missing"}, true},
	}

	check := func(t *testing.T, in ast.Internal, wantErr bool) {
		tt := &tidTable{
			byRef:       map[string]uint16{},
			aliasByName: map[string]uint16{},
			caseByName:  map[string]uint16{"target": 5},
			nodeByCase:  map[string]uint16{},
		}
		e := &enc{tt: tt, nodeNext: 100}
		got, err := emitInternalDecl(e, ref, in, span)
		if wantErr && err == nil {
			t.Fatalf("emitInternalDecl(%+v): expected error, got nil", in)
		}
		if wantErr {
			return
		}
		if err != nil {
			t.Fatalf("emitInternalDecl(%+v): %v", in, err)
		}
		if len(got) == 0 {
			t.Errorf("emitInternalDecl(%+v): expected non-empty words", in)
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { check(t, c.in, c.wantErr) })
	}
}
