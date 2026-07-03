//ff:type feature=tangl type=model
//ff:what LevelCode — a Node's crisp reporting level packed in header bits 10-9 (§5.5)
package codebook

// LevelCode is the crisp reporting level of a Node, carried in the two reserved
// header bits b10-9 of a Node packet (§5.5). It mirrors ast.Level: a `, level
// fail` / `, level review` tail modifier on a node declaration. LevelUnset is
// the zero value and encodes as 00, so every pre-v0.2 stream (whose reserved
// bits were always 0) decodes as "no explicit level" — the append-only default.
// Free-area L1 codebook value, FROZEN append-only.
type LevelCode uint16

const (
	// LevelUnset is the zero value: no explicit level was written. The engine
	// applies the default mapping (counter→Fail, except→Review); the codebook
	// only preserves the "unwritten" fact losslessly.
	LevelUnset LevelCode = 0
	// LevelFail marks the node as a reins-style FAIL contributor.
	LevelFail LevelCode = 1
	// LevelReview marks the node as a reins-style REVIEW contributor.
	LevelReview LevelCode = 2
	// 3 — reserved code, unassigned.
)
