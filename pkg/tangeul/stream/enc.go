//ff:type feature=tangl type=model
//ff:what enc — the Writer's mutable encoding state (TID table + auxiliary TID counters)
package stream

// enc is the mutable state threaded through the Writer while it serializes a
// Doc. tt holds the pass-1 declaration-order TIDs of every primary entity;
// nodeNext issues TIDs for auxiliary Func nodes (materialized function
// references, §2.1) after the primary range; textNext issues the segment-local
// ids of Text payloads. Node and Func TIDs ride in the packet bytes, but a Text
// payload's TID is stream-assigned (never encoded), so writer and reader must
// walk Text packets in the same order — hence a single monotonic textNext whose
// values the label/means/with references point at.
type enc struct {
	tt       *tidTable
	nodeNext uint16
	textNext uint16
}
