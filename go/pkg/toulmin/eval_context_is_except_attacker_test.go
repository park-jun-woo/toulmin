//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestEvalContextIsExceptAttacker — exceptions are recognised by Defeater strength or by the "except" role
package toulmin

import "testing"

func TestEvalContextIsExceptAttacker(t *testing.T) {
	ec := &evalContext{
		strMap:  map[string]Strength{"d": Defeater, "c": Defeasible, "s": Strict, "e": Defeasible},
		roleMap: map[string]string{"c": "counter", "e": "except"},
	}
	cases := map[string]bool{"d": true, "e": true, "c": false, "s": false, "unknown": false}
	for id, want := range cases {
		if got := ec.isExceptAttacker(id); got != want {
			t.Errorf("%s: got %v, want %v", id, got, want)
		}
	}
}
