//ff:func feature=tangl type=util control=sequence
//ff:what emitHeader — encode the STREAM_HEADER (version stamp + authoring locale) and SEGMENT_BOUNDARY controls
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// emitHeader encodes the two leading stream controls (§7.1): a STREAM_HEADER
// carrying the version stamp (GEUL spec version, TANGL profile id, codebook
// version — §11-1/§11-2) that the reader checks to reject an unknown major plus
// the authoring locale code (§11-11) that mode selection needs, and a
// SEGMENT_BOUNDARY (one document = one segment). The version payload words are
// packed here (not via packet.EncodeControl, which models only the framing
// controls) but use only word/codebook primitives, so bit knowledge stays out.
func emitHeader(locale uint16) []word.Word {
	geul := word.Word(uint16(codebook.GEULVersionMajor)<<8 | uint16(codebook.GEULVersionMinor))
	patch := word.Word(uint16(codebook.GEULVersionPatch) << 8)
	profile := word.Word(codebook.TANGLProfileID)
	cbook := word.Word(uint16(codebook.CodebookVersionMajor)<<8 | uint16(codebook.CodebookVersionMinor))

	header := []word.Word{word.Word(codebook.ControlStreamHeader), geul, patch, profile, cbook, word.Word(locale)}
	out := append(word.EncodeFrame(word.Frame{Class: word.ClassControl, Len: uint32(len(header))}), header...)

	boundary := []word.Word{word.Word(codebook.ControlSegmentBoundary)}
	out = append(out, word.EncodeFrame(word.Frame{Class: word.ClassControl, Len: uint32(len(boundary))})...)
	out = append(out, boundary...)
	return out
}
