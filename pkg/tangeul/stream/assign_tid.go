//ff:func feature=tangl type=util control=sequence
//ff:what assignTID — issue the next TID for a StmtRef and record it in the table
package stream

import "github.com/park-jun-woo/toulmin/pkg/tangeul/doc"

// assignTID issues the next ascending segment-local TID, records it under ref's
// canonical key, advances the counter, and returns it. It is the sole issuer of
// primary entity TIDs, keeping issue order equal to declaration order (§7.2).
func assignTID(tt *tidTable, next *uint16, ref doc.StmtRef) uint16 {
	id := *next
	*next++
	tt.byRef[tidKey(ref)] = id
	return id
}
