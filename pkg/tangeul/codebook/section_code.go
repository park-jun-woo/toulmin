//ff:type feature=tangl type=model
//ff:what SectionCode — the seven TANGEUL document section subcodes (codebook §2.4)
package codebook

// SectionCode is the SECTION subcode (the payload of a ControlSection packet),
// identifying one of the seven TANGEUL document sections in declaration order
// (TANGEUL §2.4, §7.1). It matches the StmtRef.Sec coordinate used by the doc
// layer. Free-area L1 codebook value, FROZEN append-only.
type SectionCode uint16

// The seven TANGEUL section subcodes (§2.4), assigned 1..7 (0 is unused).
const (
	// SectionSubject is `this document is x` (1).
	SectionSubject SectionCode = 1
	// SectionSee is `see a from path` declarations (2).
	SectionSee SectionCode = 2
	// SectionDefinitions is `means`/`has` definitions (3).
	SectionDefinitions SectionCode = 3
	// SectionRules is inline rule condition expressions (4).
	SectionRules SectionCode = 4
	// SectionCases is `in case of x` cases (5).
	SectionCases SectionCode = 5
	// SectionProvides is `provides x` endpoints (6).
	SectionProvides SectionCode = 6
	// SectionInternal is `on`/`every` internal triggers (7).
	SectionInternal SectionCode = 7
)
