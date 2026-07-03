//ff:type feature=tangl type=model
//ff:what ObjKind — the 2-bit discriminator of a Claim object slot (TID | NUM | SPAN)
package packet

// ObjKind classifies the object slot of a Claim block (§5.2). Because a Claim
// block itself does not carry a kind word (its four words are subject,
// predicate, object, mode in §5.2/§10 order), the 2-bit kind is stored in the
// enclosing Node header word (§5.1 reserved region, bits 8-7); the decoder reads
// it there before reading the Claim so it knows the object slot width.
type ObjKind uint16

const (
	// ObjTID is a single-word segment-local TID object (the default).
	ObjTID ObjKind = iota
	// ObjNum is a 2-word NUM32 numeric object.
	ObjNum
	// ObjSpan is a 3-word SPAN object — the descriptive-claim fallback anchor (§5.2).
	ObjSpan
)
