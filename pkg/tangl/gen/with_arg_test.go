//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWithArg — tests withArg for Spec, Const, and fallback goIdent branches
package gen

import "testing"

func TestWithArg(t *testing.T) {
	t.Run("spec branch", func(t *testing.T) {
		gc := &genContext{
			Defs: map[string]defInfo{
				"threshold": {Spec: "thresholdSpec", Const: "thresholdConst"},
			},
		}
		got := withArg(gc, "threshold")
		if got != "thresholdSpec" {
			t.Errorf("expected thresholdSpec, got %q", got)
		}
	})

	t.Run("const branch", func(t *testing.T) {
		gc := &genContext{
			Defs: map[string]defInfo{
				"threshold": {Const: "thresholdConst"},
			},
		}
		got := withArg(gc, "threshold")
		if got != "thresholdConst" {
			t.Errorf("expected thresholdConst, got %q", got)
		}
	})

	t.Run("fallback not in defs", func(t *testing.T) {
		gc := &genContext{
			Defs: map[string]defInfo{},
		}
		got := withArg(gc, "order received")
		if got != "orderReceived" {
			t.Errorf("expected orderReceived, got %q", got)
		}
	})

	t.Run("fallback empty spec and const", func(t *testing.T) {
		gc := &genContext{
			Defs: map[string]defInfo{
				"widget": {},
			},
		}
		got := withArg(gc, "widget")
		if got != "widget" {
			t.Errorf("expected widget, got %q", got)
		}
	})
}
