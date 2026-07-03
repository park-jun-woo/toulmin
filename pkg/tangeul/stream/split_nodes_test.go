//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestSplitNodes — Func nodes and primary nodes partition correctly, including the empty and single-kind cases
package stream

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
)

// TestSplitNodes drives splitNodes over an empty slice (the loop runs zero
// times), an all-Func slice, an all-primary slice and a mixed slice, checking
// that every Type==codebook.NodeFunc node lands in funcNodes (in order) and
// every other node lands in primaries (in order).
func TestSplitNodes(t *testing.T) {
	fn1 := packet.NodePacket{TID: 1, Type: codebook.NodeFunc}
	fn2 := packet.NodePacket{TID: 2, Type: codebook.NodeFunc}
	rule := packet.NodePacket{TID: 3, Type: codebook.NodeRule}
	kase := packet.NodePacket{TID: 4, Type: codebook.NodeCase}

	cases := []struct {
		name          string
		nodes         []packet.NodePacket
		wantFuncNodes []packet.NodePacket
		wantPrimaries []packet.NodePacket
	}{
		{name: "empty", nodes: nil, wantFuncNodes: nil, wantPrimaries: nil},
		{name: "all func", nodes: []packet.NodePacket{fn1, fn2}, wantFuncNodes: []packet.NodePacket{fn1, fn2}, wantPrimaries: nil},
		{name: "all primary", nodes: []packet.NodePacket{rule, kase}, wantFuncNodes: nil, wantPrimaries: []packet.NodePacket{rule, kase}},
		{name: "mixed", nodes: []packet.NodePacket{rule, fn1, kase, fn2}, wantFuncNodes: []packet.NodePacket{fn1, fn2}, wantPrimaries: []packet.NodePacket{rule, kase}},
	}
	for _, c := range cases {
		gotFuncNodes, gotPrimaries := splitNodes(c.nodes)
		if !reflect.DeepEqual(gotFuncNodes, c.wantFuncNodes) {
			t.Errorf("%s: funcNodes = %+v, want %+v", c.name, gotFuncNodes, c.wantFuncNodes)
		}
		if !reflect.DeepEqual(gotPrimaries, c.wantPrimaries) {
			t.Errorf("%s: primaries = %+v, want %+v", c.name, gotPrimaries, c.wantPrimaries)
		}
	}
}
