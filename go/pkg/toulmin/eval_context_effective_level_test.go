//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestEvalContextEffectiveLevel — tests explicit set wins over role default
package toulmin

import "testing"

func TestEvalContextEffectiveLevel(t *testing.T) {
	ec := &evalContext{
		levelMap:    map[string]Level{"a": Review, "b": Unset},
		levelSetMap: map[string]bool{"a": true, "b": true},
	}

	cases := []struct {
		name string
		id   string
		role string
		want Level
	}{
		{"explicit set returns levelMap value", "a", "counter", Review},
		{"explicit set to Unset overrides role default", "b", "counter", Unset},
		{"unset falls back to role default counter", "x", "counter", Fail},
		{"unset falls back to role default except", "x", "except", Review},
		{"unset falls back to role default other", "x", "rule", Unset},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := ec.effectiveLevel(c.id, c.role); got != c.want {
				t.Errorf("effectiveLevel(%q, %q) = %v, want %v", c.id, c.role, got, c.want)
			}
		})
	}
}
