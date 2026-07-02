//ff:func feature=policy type=engine control=sequence
//ff:what TestFormatTrace — covers empty, single, and multiple trace entries for formatTrace loop
package policy

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestFormatTrace(t *testing.T) {
	if got := formatTrace(nil); got != "" {
		t.Fatalf("empty traces: got %q, want %q", got, "")
	}

	single := []toulmin.TraceEntry{{Name: "ruleA", Activated: true}}
	if got, want := formatTrace(single), "ruleA=true"; got != want {
		t.Fatalf("single trace: got %q, want %q", got, want)
	}

	multi := []toulmin.TraceEntry{
		{Name: "ruleA", Activated: true},
		{Name: "ruleB", Activated: false},
	}
	if got, want := formatTrace(multi), "ruleA=true, ruleB=false"; got != want {
		t.Fatalf("multi trace: got %q, want %q", got, want)
	}
}
