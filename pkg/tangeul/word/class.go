//ff:type feature=tangl type=model
//ff:what Class — the 2-bit packet class carried in the frame header (Node/Edge/Payload/Control)
package word

// Class is the 2-bit packet class stored in bits 11-10 of the frame header
// word (GEUL §2.2). The four values are fixed by the grammar.
type Class uint8

const (
	// ClassNode marks a Node packet (Class = 00).
	ClassNode Class = 0
	// ClassEdge marks an Edge packet (Class = 01).
	ClassEdge Class = 1
	// ClassPayload marks a Payload packet (Class = 10).
	ClassPayload Class = 2
	// ClassControl marks a Control packet (Class = 11).
	ClassControl Class = 3
)
