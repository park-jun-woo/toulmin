//ff:func feature=tangl type=util control=sequence
//ff:what TestParseOperandSteps — parseOperandSteps resolves [count][label idx × count] and rejects truncated/out-of-range bodies
package packet

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestParseOperandSteps exercises parseOperandSteps's successful reads — a
// multi-step body resolved through labels and a zero-step body — plus its
// three guards: a missing step count, a truncated step list and an
// out-of-range step label index.
func TestParseOperandSteps(t *testing.T) {
	labels := []string{"author", "trustScore"}

	pos := 0
	steps, err := parseOperandSteps([]word.Word{2, 0, 1}, &pos, labels)
	if err != nil || pos != 3 {
		t.Fatalf("two steps: %q, pos %d, err %v", steps, pos, err)
	}
	if want := []string{"author", "trustScore"}; !reflect.DeepEqual(steps, want) {
		t.Fatalf("two steps = %q, want %q", steps, want)
	}

	pos = 0
	steps, err = parseOperandSteps([]word.Word{0}, &pos, labels)
	if err != nil || pos != 1 || len(steps) != 0 {
		t.Fatalf("zero steps: %q, pos %d, err %v", steps, pos, err)
	}

	pos = 0
	if _, err := parseOperandSteps(nil, &pos, labels); err == nil {
		t.Error("missing step count: expected error")
	}

	pos = 0
	if _, err := parseOperandSteps([]word.Word{3, 0, 1}, &pos, labels); err == nil {
		t.Error("truncated steps: expected error")
	}

	pos = 0
	if _, err := parseOperandSteps([]word.Word{1, 9}, &pos, labels); err == nil {
		t.Error("out-of-range step label: expected error")
	}
}
