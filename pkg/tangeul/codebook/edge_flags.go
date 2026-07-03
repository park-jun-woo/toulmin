//ff:type feature=tangl type=model
//ff:what EdgeFlags — the Edge flag bits (codebook §6.1 b0-b3)
package codebook

// EdgeFlags is the 4-bit flag bitset in the low bits of the Edge flag word
// (§6.1). FROZEN append-only.
type EdgeFlags uint16

// The Edge flags. Names follow the §6.1 F_* labels.
const (
	// EdgeFlagOnce (F_ONCE, b0) — `once` tick-idempotent do/undo.
	EdgeFlagOnce EdgeFlags = 1 << 0
	// EdgeFlagGate (F_GATE, b1) — a GATE16 confidence gate slot follows.
	EdgeFlagGate EdgeFlags = 1 << 1
	// EdgeFlagStrike (F_STRIKE, b2) — struck-through edge, skipped.
	EdgeFlagStrike EdgeFlags = 1 << 2
	// EdgeFlagOrdered (F_ORDERED, b3) — belongs to a `1.` ordered sequence.
	EdgeFlagOrdered EdgeFlags = 1 << 3
)
