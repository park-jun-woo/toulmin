//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestBoolRuleFunc — tests boolRuleFunc maps true → +1 and false → -1 and passes evidence through
package toulmin

import "testing"

func TestBoolRuleFunc(t *testing.T) {
	cases := []struct {
		in   bool
		want float64
	}{
		{true, 1.0},
		{false, -1.0},
	}
	for _, c := range cases {
		in := c.in
		f := boolRuleFunc(func(ctx Context, specs Specs) (bool, any) { return in, in })
		v, evidence := f(NewContext(), nil)
		if v != c.want {
			t.Errorf("%v: expected %v, got %v", c.in, c.want, v)
		}
		if evidence != c.in {
			t.Errorf("%v: expected evidence passed through, got %v", c.in, evidence)
		}
	}
}
