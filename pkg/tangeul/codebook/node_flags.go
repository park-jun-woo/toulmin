//ff:type feature=tangl type=model
//ff:what NodeFlags — the Node header flag bits (codebook §5.1 b0-b5 + F_WITH b6)
package codebook

// NodeFlags is the bitset in the low bits of the Node header word (§5.1). Bits
// b0-b5 are the six flags defined by §5.1; each set flag introduces its
// conditional slot in flag order. F_WITH (b6) is allocated from the reserved
// region for the `with` clause slot (§11-4 basis, consumed in Phase 003).
// FROZEN append-only.
type NodeFlags uint16

// The Node header flags. Names follow the §5.1 F_* labels.
const (
	// NodeFlagFunc (F_FUNC, b0) — a FuncRef warrant-function slot follows.
	NodeFlagFunc NodeFlags = 1 << 0
	// NodeFlagQual (F_QUAL, b1) — a VFIX16 qualifier slot follows.
	NodeFlagQual NodeFlags = 1 << 1
	// NodeFlagCheck (F_CHECK, b2) — a TID `checking` target Case slot follows.
	NodeFlagCheck NodeFlags = 1 << 2
	// NodeFlagClaim (F_CLAIM, b3) — a Claim block (§5.2) follows.
	NodeFlagClaim NodeFlags = 1 << 3
	// NodeFlagVerdict (F_VERDICT, b4) — a VFIX16 initial verdict slot follows.
	NodeFlagVerdict NodeFlags = 1 << 4
	// NodeFlagStrike (F_STRIKE, b5) — struck-through annotation, skipped.
	NodeFlagStrike NodeFlags = 1 << 5
	// NodeFlagWith (F_WITH, b6) — a `with` clause slot (reserved, Phase 003).
	NodeFlagWith NodeFlags = 1 << 6
)
