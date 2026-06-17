//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestNodeEventTypeString — String returns the name for all three event types
package toulmin

import "testing"

func TestNodeEventTypeString(t *testing.T) {
	cases := []struct {
		in   NodeEventType
		want string
	}{
		{Active, "Active"},
		{Defeated, "Defeated"},
		{Inactive, "Inactive"},
	}
	for _, c := range cases {
		if got := c.in.String(); got != c.want {
			t.Errorf("String() = %q, want %q", got, c.want)
		}
	}
}
