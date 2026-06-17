//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestClassifyEvent — classifyEvent maps active flag and verdict to all three event types
package toulmin

import "testing"

func TestClassifyEvent(t *testing.T) {
	cases := []struct {
		active  bool
		verdict float64
		want    NodeEventType
	}{
		{false, 1.0, Inactive},
		{false, 0.0, Inactive},
		{true, 0.5, Active},
		{true, 0.0, Defeated},
		{true, -0.5, Defeated},
	}
	for _, c := range cases {
		if got := classifyEvent(c.active, c.verdict); got != c.want {
			t.Errorf("classifyEvent(%v, %v) = %v, want %v", c.active, c.verdict, got, c.want)
		}
	}
}
