//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestGateRoundTrip — decode(encode(gate)) equals gate across all operators and sample percents
package word

import "testing"

func TestGateRoundTrip(t *testing.T) {
	cases := []Gate{
		{GateAtLeast, 0}, {GateAtLeast, 75}, {GateAtLeast, 100},
		{GateAbove, 1}, {GateAbove, 50},
		{GateLessThan, 33}, {GateLessThan, 100},
		{GateAtMost, 0}, {GateAtMost, 99},
	}
	for _, want := range cases {
		w, err := EncodeGate(want)
		got, derr := DecodeGate(w)
		if err != nil || derr != nil || got != want {
			t.Fatalf("roundtrip %+v: got %+v enc=%v dec=%v", want, got, err, derr)
		}
	}
}
