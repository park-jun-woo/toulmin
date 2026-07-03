//ff:type feature=tangl type=model
//ff:what GateOp — the 2-bit threshold operator of a confidence GATE16
package word

// GateOp is the 2-bit threshold operator of a GATE16 (GEUL §3.2), matching
// TANGEUL's THRESHOLD_OP. Mapping between these values and the ast.Certainty.Op
// strings is the codebook's concern, not this package's.
type GateOp uint8

const (
	// GateAtLeast is "at least N%" (op = 0).
	GateAtLeast GateOp = 0
	// GateAbove is "above N%" (op = 1).
	GateAbove GateOp = 1
	// GateLessThan is "less than N%" (op = 2).
	GateLessThan GateOp = 2
	// GateAtMost is "at most N%" (op = 3).
	GateAtMost GateOp = 3
)
