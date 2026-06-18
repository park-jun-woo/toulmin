//ff:func feature=engine type=engine control=sequence
//ff:what TestRunHandlerError — handler error or panic stops Run and propagates with node name
package toulmin

import (
	"fmt"
	"strings"
	"testing"
)

func TestRunHandlerError(t *testing.T) {
	secondFired := false
	g := NewGraph("err")
	g.Rule(WarrantA).RunOn(func(ctx Context, self TraceEntry, trace []TraceEntry) error {
		return fmt.Errorf("boom")
	})
	g.Rule(RebuttalB).RunOn(func(ctx Context, self TraceEntry, trace []TraceEntry) error {
		secondFired = true
		return nil
	})
	_, trace, err := g.Run(NewContext())
	if err == nil {
		t.Fatal("expected error from handler")
	}
	if !strings.Contains(err.Error(), "WarrantA") {
		t.Errorf("error should name the node: %v", err)
	}
	if secondFired {
		t.Error("Run must stop before firing later handlers")
	}
	if trace == nil {
		t.Error("on handler error Run must still return the pre-dispatch trace")
	}

	g2 := NewGraph("panic")
	g2.Rule(WarrantA).RunOn(func(ctx Context, self TraceEntry, trace []TraceEntry) error {
		panic("kaboom")
	})
	_, trace2, err2 := g2.Run(NewContext())
	if err2 == nil {
		t.Fatal("expected error from panicking handler")
	}
	if !strings.Contains(err2.Error(), "panicked") {
		t.Errorf("panic should convert to error: %v", err2)
	}
	if trace2 == nil {
		t.Error("on panic Run must still return the pre-dispatch trace")
	}
}
