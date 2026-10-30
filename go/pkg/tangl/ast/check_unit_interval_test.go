//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestCheckUnitInterval — accepts [0, 1] at 4 decimals; rejects NaN, out-of-range, and finer values
package ast

import (
	"math"
	"strings"
	"testing"
)

func TestCheckUnitInterval(t *testing.T) {
	for _, v := range []float64{0, 0.0001, 0.5, 0.8, 0.1235, 1} {
		if err := CheckUnitInterval(v); err != nil {
			t.Errorf("CheckUnitInterval(%v) = %v, want nil", v, err)
		}
	}
	cases := []struct {
		v    float64
		want string
	}{
		{math.NaN(), "NaN"},
		{-0.1, "out of range"},
		{1.0001, "out of range"},
		{0.12345, "more than 4 decimal places"},
	}
	for _, tc := range cases {
		err := CheckUnitInterval(tc.v)
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("CheckUnitInterval(%v) = %v, want error containing %q", tc.v, err, tc.want)
		}
	}
}
