//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestRegisterFuncNodes — records alias/name per Func node TID, zero and non-zero alias, empty input
package stream

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
)

// TestRegisterFuncNodes drives registerFuncNodes directly (white-box) over the
// loop's zero-iteration case (no Func nodes: d.funcRef stays empty) and its
// nonzero-iteration case covering both the fn.Func == nil branch (alias stays
// 0, a local reference) and the fn.Func != nil branch (alias is read from the
// F_FUNC slot), plus a TID with no recorded label (firstLabel's empty-string
// path).
func TestRegisterFuncNodes(t *testing.T) {
	tests := []struct {
		name      string
		funcNodes []packet.NodePacket
		labels    map[uint16][]string
		want      map[uint16]funcInfo
	}{
		{
			name:      "no func nodes",
			funcNodes: nil,
			labels:    map[uint16][]string{1: {"withdraw"}},
			want:      map[uint16]funcInfo{},
		},
		{
			name: "nil Func slot and labeled node",
			funcNodes: []packet.NodePacket{
				{TID: 5, Func: nil},
			},
			labels: map[uint16][]string{5: {"withdraw", "extra"}},
			want: map[uint16]funcInfo{
				5: {alias: 0, name: "withdraw"},
			},
		},
		{
			name: "non-nil Func slot with alias",
			funcNodes: []packet.NodePacket{
				{TID: 7, Func: &packet.FuncRef{Alias: 3}},
			},
			labels: map[uint16][]string{7: {"log"}},
			want: map[uint16]funcInfo{
				7: {alias: 3, name: "log"},
			},
		},
		{
			name: "unlabeled TID falls back to empty name",
			funcNodes: []packet.NodePacket{
				{TID: 9, Func: &packet.FuncRef{Alias: 2}},
			},
			labels: map[uint16][]string{},
			want: map[uint16]funcInfo{
				9: {alias: 2, name: ""},
			},
		},
		{
			name: "multiple func nodes in one statement",
			funcNodes: []packet.NodePacket{
				{TID: 10, Func: nil},
				{TID: 11, Func: &packet.FuncRef{Alias: 4}},
			},
			labels: map[uint16][]string{10: {"a"}, 11: {"b"}},
			want: map[uint16]funcInfo{
				10: {alias: 0, name: "a"},
				11: {alias: 4, name: "b"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newDec()
			registerFuncNodes(d, tt.funcNodes, tt.labels)
			if !reflect.DeepEqual(d.funcRef, tt.want) {
				t.Errorf("funcRef = %+v, want %+v", d.funcRef, tt.want)
			}
		})
	}
}
