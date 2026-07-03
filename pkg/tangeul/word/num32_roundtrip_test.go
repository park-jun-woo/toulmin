//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestNum32RoundTrip — decode(encode(n)) equals n across the 16-bit boundary and max uint32
package word

import "testing"

func TestNum32RoundTrip(t *testing.T) {
	for _, n := range []uint32{0, 1, 1023, 65535, 65536, 0xFFFFFFFF} {
		ws := EncodeNum32(n)
		got, err := DecodeNum32(ws[:])
		if err != nil || got != n {
			t.Fatalf("roundtrip %d: got %d err=%v", n, got, err)
		}
	}
}
