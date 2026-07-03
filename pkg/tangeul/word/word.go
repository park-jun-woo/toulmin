//ff:type feature=tangl type=model
//ff:what Word — a single 16-bit GEUL word (the atomic unit of a TANGEUL stream)
package word

// Word is one 16-bit GEUL word. TANGEUL packets are sequences of Words encoded
// Big Endian, MSB first, word aligned (GEUL §2.1). All bit-layout knowledge in
// the TANGEUL implementation is confined to this package (and codebook), so the
// packet/stream/doc layers deal only with typed fields and stay immune to the
// L0 bit-placement churn documented as Proposal in the grammar (§11-1).
type Word uint16
