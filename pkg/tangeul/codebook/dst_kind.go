//ff:type feature=tangl type=model
//ff:what DstKind — the slot kind an Edge dst occupies, as decided by its PropCode
package codebook

// DstKind classifies the slot that an Edge packet's dst occupies. The PropCode
// alone decides this (§2.3): `do`/`undo` carry a FuncRef (2 words), `label` and
// `importedFrom` carry Text, `provenance` carries a Block, and every other
// predicate carries a single-word TID. A decoder consults PropDstKind to know
// how to read word 4 onward of an Edge.
type DstKind int

const (
	// DstTID is a single-word segment-local TID reference (the default).
	DstTID DstKind = iota
	// DstFuncRef is a 2-word FuncRef (Alias TID + VerbCode) — do/undo only.
	DstFuncRef
	// DstText is a Text reference — label, importedFrom (path).
	DstText
	// DstBlock is a structured block — provenance.
	DstBlock
)
