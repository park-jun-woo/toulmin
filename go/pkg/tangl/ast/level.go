//ff:type feature=tangl type=model
//ff:what Level — a node's crisp reporting level (unset/fail/review)
package ast

// Level classifies a Node's crisp reporting level for reins-style
// aggregation. It aligns numerically and semantically with toulmin.Level
// but is a distinct tool-chain type; the ast package never imports toulmin.
type Level int

const (
	// LevelUnset is the zero value: no explicit level was written.
	LevelUnset Level = iota
	// LevelFail maps to toulmin.Fail at codegen time.
	LevelFail
	// LevelReview maps to toulmin.Review at codegen time.
	LevelReview
)
