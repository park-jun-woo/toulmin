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
	g.Rule(WarrantA).OnActive(func(ctx Context, ev NodeEvent, view RunView) error {
		return fmt.Errorf("boom")
	})
	g.Counter(RebuttalB).OnActive(func(ctx Context, ev NodeEvent, view RunView) error {
		secondFired = true
		return nil
	})
	_, view, err := g.Run(NewContext())
	if err == nil {
		t.Fatal("expected error from handler")
	}
	if !strings.Contains(err.Error(), "WarrantA") {
		t.Errorf("error should name the node: %v", err)
	}
	if !strings.Contains(err.Error(), "Active") {
		t.Errorf("error should name the event type: %v", err)
	}
	if secondFired {
		t.Error("Run must stop before firing later handlers")
	}
	if view == nil {
		t.Error("on handler error Run must still return the pre-dispatch RunView")
	}

	g2 := NewGraph("panic")
	g2.Rule(WarrantA).OnActive(func(ctx Context, ev NodeEvent, view RunView) error {
		panic("kaboom")
	})
	_, view2, err2 := g2.Run(NewContext())
	if err2 == nil {
		t.Fatal("expected error from panicking handler")
	}
	if !strings.Contains(err2.Error(), "panicked") {
		t.Errorf("panic should convert to error: %v", err2)
	}
	if view2 == nil {
		t.Error("on panic Run must still return the pre-dispatch RunView")
	}
}
