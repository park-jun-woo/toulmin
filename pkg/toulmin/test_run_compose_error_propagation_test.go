//ff:func feature=engine type=engine control=sequence
//ff:what TestRunComposeErrorPropagation — a sub-graph handler error propagates wrapped as `run ... → ...`
package toulmin

import (
	"fmt"
	"strings"
	"testing"
)

func TestRunComposeErrorPropagation(t *testing.T) {
	subRule := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	sub := NewGraph("sub")
	sub.Rule(subRule).RunOn(func(t Trace) error {
		return fmt.Errorf("sub boom")
	})

	active := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	parent := NewGraph("parent")
	parent.Rule(active).Run(sub)

	results, trace, err := parent.Run(NewContext())
	if err == nil {
		t.Fatal("expected sub-graph error to propagate")
	}
	if !strings.Contains(err.Error(), "run ") || !strings.Contains(err.Error(), "→") {
		t.Errorf("error should be wrapped as `run ... → ...`: %v", err)
	}
	if !strings.Contains(err.Error(), "sub boom") {
		t.Errorf("error should wrap the underlying cause: %v", err)
	}
	if results == nil || trace.All() == nil {
		t.Error("on sub-Run error the parent must still return its pre-dispatch results and trace")
	}
}
