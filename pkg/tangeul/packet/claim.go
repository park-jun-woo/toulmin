//ff:type feature=tangl type=model
//ff:what Claim — a node's proposition block (subject/predicate/object/mode, §5.2)
package packet

import "github.com/park-jun-woo/toulmin/pkg/tangeul/word"

// Claim is the proposition a node carries — its cache key (§5.2). It encodes as
// four logical slots in §5.2/§10 order: subject TID, predicate CODE, object, and
// mode CODE. The object is one of TID / NUM32 / SPAN, selected by Kind; Kind is
// not emitted inside the block but is placed in the Node header word by the node
// encoder (see ObjKind). The type-ladder fallback (§5.2, §5.6) is a descriptive
// claim with Subject=Predicate=0, Kind=ObjSpan and ObjSpan anchoring the label.
type Claim struct {
	Subject   TID
	Predicate uint16
	Kind      ObjKind
	ObjTID    TID
	ObjNum    uint32
	ObjSpan   word.Span
	Mode      uint16
}
