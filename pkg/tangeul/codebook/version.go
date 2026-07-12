//ff:type feature=tangl type=model
//ff:what version stamp — GEUL version, TANGEUL profile ID and codebook version constants
package codebook

// The stream version stamp written into the STREAM_HEADER Control packet
// (§2.4, §7.1) and checked by the reader, which rejects an unknown major
// version. Three parts: the GEUL primitive spec version this build targets, the
// TANGEUL Free-area profile ID, and this codebook's own version.
//
// All values are Proposal placeholders per §11-1/§11-2; the codebook version is
// v0.3 and is FROZEN append-only from here.
//
// Minor bumped 2→3 for the Phase003 judgment-expression extension (EXPR record
// kind widened to 3 bits with the existing bit patterns preserved, plus the new
// QuantOp code) — a minor, not major, because every existing stream's bit
// interpretation is unchanged and the change is purely additive (checkVersion
// rejects only a major mismatch, so old readers keep loading old streams).
const (
	// GEULVersionMajor is the targeted GEUL spec major version (GEUL v0.3.0).
	GEULVersionMajor uint8 = 0
	// GEULVersionMinor is the targeted GEUL spec minor version (GEUL v0.3.0).
	GEULVersionMinor uint8 = 3
	// GEULVersionPatch is the targeted GEUL spec patch version (GEUL v0.3.0).
	GEULVersionPatch uint8 = 0

	// TANGLProfileID identifies the TANGEUL Free-area profile (Proposal).
	TANGLProfileID uint16 = 0x0001

	// CodebookVersionMajor is this codebook's major version (v0.2).
	CodebookVersionMajor uint8 = 0
	// CodebookVersionMinor is this codebook's minor version (v0.3).
	CodebookVersionMinor uint8 = 3
)

// Authoring locale codes stamped into the STREAM_HEADER (§11-11) so a mode-2
// decoder can tell the authoring locale (verbatim, byte-identical) from a
// cross locale (normalized rendering). 0 means unknown/unset.
const (
	// LocaleUnknown marks a stream with no recorded authoring locale.
	LocaleUnknown uint16 = 0
	// LocaleEN is the English authoring locale code.
	LocaleEN uint16 = 1
	// LocaleKO is the Korean authoring locale code.
	LocaleKO uint16 = 2
)
