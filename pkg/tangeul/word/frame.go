//ff:type feature=tangl type=model
//ff:what Frame — a decoded frame header (Class + body word count, 32-bit extension aware)
package word

// Frame is a decoded TANGEUL frame header (GEUL §2.3). Word 1 packs the Free
// prefix (0001), the 2-bit Class and a 10-bit Len. Len is the body word count;
// values below the 1023 sentinel fit the single header word, and a body of
// 1023 words or more triggers the 32-bit extension (sentinel Len + NUM32).
// Len is kept as uint32 so callers never see the sentinel encoding.
type Frame struct {
	Class Class
	Len   uint32
}

// framePrefix is the Free-area 4-bit prefix (0001) occupying bits 15-12 of the
// frame header word (GEUL §2.2).
const framePrefix = 0x1

// frameLenSentinel is the 10-bit Len value that signals a 32-bit body-length
// extension follows as a NUM32 (GEUL §2.3).
const frameLenSentinel = 1023
