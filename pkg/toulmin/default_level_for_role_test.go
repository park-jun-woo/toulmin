//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestDefaultLevelForRole — tests role→default Level mapping (counter/except/other)
package toulmin

import "testing"

func TestDefaultLevelForRole(t *testing.T) {
	cases := []struct {
		role string
		want Level
	}{
		{"counter", Fail},
		{"except", Review},
		{"rule", Unset},
		{"", Unset},
		{"unknown", Unset},
	}
	for _, c := range cases {
		if got := defaultLevelForRole(c.role); got != c.want {
			t.Errorf("defaultLevelForRole(%q) = %v, want %v", c.role, got, c.want)
		}
	}
}
