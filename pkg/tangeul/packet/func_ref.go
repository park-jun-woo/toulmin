//ff:type feature=tangl type=model
//ff:what FuncRef — a 2-word function reference (Alias TID + VerbCode) with its source label
package packet

// FuncRef is the two-word function slot of a `do`/`undo`/`using` (§4): word A is
// the package Alias TID (0x0000 for a local/self-subject function) and word B is
// the VerbCode. Verb is provisionally unbound (0x0000 allowed, README §6) — the
// packet layer never interprets it. Label holds the function's source spelling;
// it is serialized separately as a label Text packet by the stream layer and is
// NOT emitted by encodeFuncRef, so a packet round-trip preserves only Alias/Verb.
type FuncRef struct {
	Alias TID
	Verb  uint16
	Label string
}
