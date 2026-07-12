//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitOperandSteps — golden [kind][step count][label idx × step] layout over empty and multi-step operands
package packet

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestEmitOperandSteps locks the stepped-operand sub-record layout (§design 2):
// the OperandKind word, the step count, then one interned label index per step
// — a zero-step operand emitting just the kind and count, repeated spellings
// sharing one label entry across a packet-local table.
func TestEmitOperandSteps(t *testing.T) {
	cases := []struct {
		kind  OperandKind
		steps []string
		want  []word.Word
	}{
		{OperandPath, nil, []word.Word{0, 0}},
		{OperandPath, []string{"author", "trustScore"}, []word.Word{0, 2, 0, 1}},
		{OperandSpec, []string{"author", "min price"}, []word.Word{1, 2, 0, 2}},
	}
	var labels []string
	for _, c := range cases {
		got := emitOperandSteps(c.kind, c.steps, &labels)
		if !reflect.DeepEqual(got, c.want) {
			t.Fatalf("emitOperandSteps(%d, %q) = %v, want %v", c.kind, c.steps, got, c.want)
		}
	}
	if want := []string{"author", "trustScore", "min price"}; !reflect.DeepEqual(labels, want) {
		t.Fatalf("labels %q, want %q", labels, want)
	}
}
