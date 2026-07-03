//ff:type feature=tangl type=model
//ff:what PropCode — the seventeen TANGL Edge-packet predicate codes (codebook §2.3)
package codebook

// PropCode is the predicate code carried in word 2 of an Edge packet
// (TANGEUL §2.3, §6.1). It also selects the dst slot kind (see PropDstKind).
// Free-area L1 codebook value, FROZEN append-only.
type PropCode uint16

// The seventeen TANGL PropCode codes (§2.3), assigned 0x0101..0x0111.
const (
	// PropAttacks is `don't X when Y` — RuleNode → RuleNode (0x0101).
	PropAttacks PropCode = 0x0101
	// PropDo is `do fn when n` — RuleNode → FuncRef (0x0102).
	PropDo PropCode = 0x0102
	// PropUndo is `undo fn when n` — RuleNode → FuncRef (0x0103).
	PropUndo PropCode = 0x0103
	// PropRun is `run c when n` — RuleNode → Case (0x0104).
	PropRun PropCode = 0x0104
	// PropChecks is `... checking c` — RuleNode → Case, pure (0x0105).
	PropChecks PropCode = 0x0105
	// PropRequires is `f is required` — Case → Field (0x0106).
	PropRequires PropCode = 0x0106
	// PropMeans is a Definitions `means` — Term → value/Type (0x0107).
	PropMeans PropCode = 0x0107
	// PropHas is a Definitions `has` — Term → Field (0x0108).
	PropHas PropCode = 0x0108
	// PropRuns is an Endpoint `run c` — Endpoint → Case, Run (0x0109).
	PropRuns PropCode = 0x0109
	// PropEvaluates is an Endpoint `check c` — Endpoint → Case, Evaluate (0x010A).
	PropEvaluates PropCode = 0x010A
	// PropTriggers is a Trigger `run c` — Trigger → Case (0x010B).
	PropTriggers PropCode = 0x010B
	// PropUntil is `every … until c` — Trigger → Case (0x010C).
	PropUntil PropCode = 0x010C
	// PropSameAs binds a global ID — local TID → global QID/Synset (0x010D).
	PropSameAs PropCode = 0x010D
	// PropImportedFrom is `see a from path` — Alias → path(Text) (0x010E).
	PropImportedFrom PropCode = 0x010E
	// PropLabel is a string label — Node → Text (0x010F).
	PropLabel PropCode = 0x010F
	// PropProvenance is rule origin/premises meta — Node → provenance block (0x0110).
	PropProvenance PropCode = 0x0110
	// PropSupersedes is `` `X` supersedes `Y` `` — RuleNode → RuleNode (0x0111).
	PropSupersedes PropCode = 0x0111
)
