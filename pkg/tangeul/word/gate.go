//ff:type feature=tangl type=model
//ff:what Gate — a decoded confidence gate (operator + integer percent)
package word

// Gate is a decoded confidence gate (GATE16, GEUL §3.2). The word packs the
// 2-bit operator in bits 15-14, seven reserved bits (must be 0), and a 7-bit
// integer percent (0-100) in bits 6-0. Example: "if at least 75% certain" is
// Op=GateAtLeast, Percent=75, encoding to 0x004B.
type Gate struct {
	Op      GateOp
	Percent uint8
}

// gateMaxPercent is the highest legal confidence percent in a Gate.
const gateMaxPercent = 100
