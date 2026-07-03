//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestVFixRoundTrip — Float(VFixFromFloat(f)) equals f for values representable at 4 decimals
package word

import "testing"

func TestVFixRoundTrip(t *testing.T) {
	for _, f := range []float64{-1.0, -0.5, 0.0, 0.02, 0.8, 1.0} {
		v, err := VFixFromFloat(f)
		if err != nil || v.Float() != f {
			t.Fatalf("roundtrip %v: got %v err=%v", f, v.Float(), err)
		}
	}
}
