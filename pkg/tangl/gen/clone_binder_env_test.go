//ff:func feature=tangl type=codegen control=sequence
//ff:what TestCloneBinderEnv — tests binder-env copy-and-add without mutating the caller's env
package gen

import "testing"

func TestCloneBinderEnv(t *testing.T) {
	outer := map[string]binderInfo{"a": {Elem: "e0", Idx: "i0"}}
	inner := cloneBinderEnv(outer, "b", binderInfo{Elem: "e1", Idx: "i1"})
	if len(inner) != 2 || inner["a"].Elem != "e0" || inner["b"].Elem != "e1" {
		t.Errorf("clone should carry the outer binder and add the new one: %+v", inner)
	}
	if len(outer) != 1 {
		t.Errorf("caller's env must not be mutated: %+v", outer)
	}
	shadow := cloneBinderEnv(outer, "a", binderInfo{Elem: "e2", Idx: "i2"})
	if shadow["a"].Elem != "e2" {
		t.Errorf("re-binding the same name should shadow: %+v", shadow)
	}
}
