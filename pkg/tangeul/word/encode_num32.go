//ff:func feature=tangl type=util control=sequence
//ff:what EncodeNum32 — encode a uint32 as a 2-word NUM32 (high word first)
package word

// EncodeNum32 encodes n as a NUM32: two words, high 16 bits first, low 16 bits
// second (big-endian word order, GEUL §2.4).
func EncodeNum32(n uint32) [num32Words]Word {
	return [num32Words]Word{
		Word(n >> 16),
		Word(n & 0xFFFF),
	}
}
