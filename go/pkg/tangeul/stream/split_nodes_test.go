//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestSplitNodes — Func nodes and primary nodes partition correctly, including the empty and single-kind cases
package stream

import (
	"reflect"
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// TestSplitNodes drives splitNodes over an empty slice (the loop runs zero
// times), an all-Func slice, an all-primary slice and a mixed slice, checking
// that every Type==toulmin.NodeFunc node lands in funcNodes (in order) and
// every other node lands in primaries (in order).
func TestSplitNodes(t *testing.T) {
	fn1 := toulmin.NodePacket{TID: 1, Type: toulmin.NodeFunc}
	fn2 := toulmin.NodePacket{TID: 2, Type: toulmin.NodeFunc}
	rule := toulmin.NodePacket{TID: 3, Type: toulmin.NodeRule}
	kase := toulmin.NodePacket{TID: 4, Type: toulmin.NodeCase}

	cases := []struct {
		name          string
		nodes         []toulmin.NodePacket
		wantFuncNodes []toulmin.NodePacket
		wantPrimaries []toulmin.NodePacket
	}{
		{name: "empty", nodes: nil, wantFuncNodes: nil, wantPrimaries: nil},
		{name: "all func", nodes: []toulmin.NodePacket{fn1, fn2}, wantFuncNodes: []toulmin.NodePacket{fn1, fn2}, wantPrimaries: nil},
		{name: "all primary", nodes: []toulmin.NodePacket{rule, kase}, wantFuncNodes: nil, wantPrimaries: []toulmin.NodePacket{rule, kase}},
		{name: "mixed", nodes: []toulmin.NodePacket{rule, fn1, kase, fn2}, wantFuncNodes: []toulmin.NodePacket{fn1, fn2}, wantPrimaries: []toulmin.NodePacket{rule, kase}},
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
