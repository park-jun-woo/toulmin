//ff:func feature=engine type=engine control=sequence
//ff:what TestSelectHandler — selectHandler returns the handler matching each event type
package toulmin

import "testing"

func TestSelectHandler(t *testing.T) {
	active := func(ctx Context, ev NodeEvent, view RunView) error { return nil }
	defeated := func(ctx Context, ev NodeEvent, view RunView) error { return nil }
	inactive := func(ctx Context, ev NodeEvent, view RunView) error { return nil }
	r := &RuleMeta{OnActive: active, OnDefeated: defeated, OnInactive: inactive}

	if h := selectHandler(r, Active); h == nil {
		t.Error("Active should select OnActive")
	}
	if h := selectHandler(r, Defeated); h == nil {
		t.Error("Defeated should select OnDefeated")
	}
	if h := selectHandler(r, Inactive); h == nil {
		t.Error("Inactive should select OnInactive")
	}

	empty := &RuleMeta{}
	if h := selectHandler(empty, Active); h != nil {
		t.Error("nil handler should be returned when none registered")
	}
}
