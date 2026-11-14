//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestValidInteger — table-driven check of validInteger across integer and non-integer kinds
package types

import "testing"

func TestValidInteger(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want bool
	}{
		{"int", int(1), true},
		{"int8", int8(1), true},
		{"int16", int16(1), true},
		{"int32", int32(1), true},
		{"int64", int64(1), true},
		{"uint", uint(1), true},
		{"uint8", uint8(1), true},
		{"uint16", uint16(1), true},
		{"uint32", uint32(1), true},
		{"uint64", uint64(1), true},
		{"float", 1.5, false},
		{"string", "1", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := validInteger(c.v)
			if got != c.want {
				t.Errorf("validInteger(%#v) = %v, want %v", c.v, got, c.want)
			}
		})
	}
}
