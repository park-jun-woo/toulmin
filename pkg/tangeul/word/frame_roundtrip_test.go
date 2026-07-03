//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestFrameRoundTrip — decode(encode(frame)) equals frame across the extension boundary
package word

import "testing"

func TestFrameRoundTrip(t *testing.T) {
	cases := []Frame{
		{Class: ClassNode, Len: 0},
		{Class: ClassEdge, Len: 1},
		{Class: ClassPayload, Len: 1022},
		{Class: ClassControl, Len: 1023},   // extension boundary
		{Class: ClassPayload, Len: 100000}, // large extension
	}
	for _, want := range cases {
		ws := EncodeFrame(want)
		got, n, err := DecodeFrame(ws)
		if err != nil || n != len(ws) || got != want {
			t.Fatalf("roundtrip %+v: got %+v n=%d (words=%d) err=%v", want, got, n, len(ws), err)
		}
	}
}
