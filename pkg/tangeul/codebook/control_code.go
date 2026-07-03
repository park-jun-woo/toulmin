//ff:type feature=tangl type=model
//ff:what ControlCode — the seven TANGEUL Control-packet subcodes (codebook §2.4)
package codebook

// ControlCode is the subcode of a Control packet (Class 11), carried in word 2
// (TANGEUL §2.4, §7.1). Free-area L1 codebook value, FROZEN append-only.
type ControlCode uint16

// The seven TANGEUL Control subcodes (§2.4), assigned 0x00..0x06.
const (
	// ControlStreamHeader carries GEUL version, codebook decl and TANGEUL profile ID (0x00).
	ControlStreamHeader ControlCode = 0x00
	// ControlSegmentBoundary is a resync point; one document = one segment (0x01).
	ControlSegmentBoundary ControlCode = 0x01
	// ControlCodebookDecl declares the codebook in effect (0x02).
	ControlCodebookDecl ControlCode = 0x02
	// ControlSection opens a section; the subcode is a SectionCode (0x03).
	ControlSection ControlCode = 0x03
	// ControlCaseBegin opens a case, carrying the Case TID (0x04).
	ControlCaseBegin ControlCode = 0x04
	// ControlCaseEnd closes a case, carrying the Case TID (0x05).
	ControlCaseEnd ControlCode = 0x05
	// ControlProse marks a prose span (body = SPAN, §7.4) (0x06).
	ControlProse ControlCode = 0x06
)
