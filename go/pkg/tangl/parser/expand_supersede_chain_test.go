//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestExpandSupersedeChain — every earlier→later pair in order; only the first edge carries the chain copy
package parser

import (
	"reflect"
	"testing"
)

func TestExpandSupersedeChain(t *testing.T) {
	names := []string{"A", "B", "C", "D"}
	edges := expandSupersedeChain(names, 7)
	want := [][2]string{{"A", "B"}, {"A", "C"}, {"A", "D"}, {"B", "C"}, {"B", "D"}, {"C", "D"}}
	if len(edges) != len(want) {
		t.Fatalf("len = %d, want %d", len(edges), len(want))
	}
	for i, pair := range want {
		if edges[i].Superseder != pair[0] || edges[i].Superseded != pair[1] || edges[i].Line != 7 {
			t.Fatalf("edge[%d] = %+v, want %v", i, edges[i], pair)
		}
		if i > 0 && edges[i].Chain != nil {
			t.Fatalf("edge[%d] must not carry the chain", i)
		}
	}
	if !reflect.DeepEqual(edges[0].Chain, names) {
		t.Fatalf("chain = %v", edges[0].Chain)
	}
	names[0] = "Z"
	if edges[0].Chain[0] != "A" {
		t.Fatal("chain must be a copy of the input names")
	}
	if got := expandSupersedeChain([]string{"solo"}, 1); len(got) != 0 {
		t.Fatalf("single name = %+v, want no edges", got)
	}
}
