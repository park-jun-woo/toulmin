//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestRegisterFuncNodes — records alias/name per Func node TID, zero and non-zero alias, empty input
package stream

import (
	"reflect"
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"
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
		funcNodes []toulmin.NodePacket
		labels    map[toulmin.TID][]string
		want      map[toulmin.TID]funcInfo
	}{
		{
			name:      "no func nodes",
			funcNodes: nil,
			labels:    map[toulmin.TID][]string{1: {"withdraw"}},
			want:      map[toulmin.TID]funcInfo{},
		},
		{
			name: "nil Func slot and labeled node",
			funcNodes: []toulmin.NodePacket{
				{TID: 5, Func: nil},
			},
			labels: map[toulmin.TID][]string{5: {"withdraw", "extra"}},
			want: map[toulmin.TID]funcInfo{
				5: {alias: 0, name: "withdraw"},
			},
		},
		{
			name: "non-nil Func slot with alias",
			funcNodes: []toulmin.NodePacket{
				{TID: 7, Func: &toulmin.FuncRef{Alias: 3}},
			},
			labels: map[toulmin.TID][]string{7: {"log"}},
			want: map[toulmin.TID]funcInfo{
				7: {alias: 3, name: "log"},
			},
		},
		{
			name: "unlabeled TID falls back to empty name",
			funcNodes: []toulmin.NodePacket{
				{TID: 9, Func: &toulmin.FuncRef{Alias: 2}},
			},
			labels: map[toulmin.TID][]string{},
			want: map[toulmin.TID]funcInfo{
				9: {alias: 2, name: ""},
			},
		},
		{
			name: "multiple func nodes in one statement",
			funcNodes: []toulmin.NodePacket{
				{TID: 10, Func: nil},
				{TID: 11, Func: &toulmin.FuncRef{Alias: 4}},
			},
			labels: map[toulmin.TID][]string{10: {"a"}, 11: {"b"}},
			want: map[toulmin.TID]funcInfo{
				10: {alias: 0, name: "a"},
				11: {alias: 4, name: "b"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newDec(testWire)
			registerFuncNodes(d, tt.funcNodes, tt.labels)
			if !reflect.DeepEqual(d.funcRef, tt.want) {
				t.Errorf("funcRef = %+v, want %+v", d.funcRef, tt.want)
			}
		})
	}
}
