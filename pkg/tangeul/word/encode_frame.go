//ff:func feature=tangl type=util control=sequence
//ff:what EncodeFrame — encode a frame header, emitting the 32-bit extension when the body is large
package word

// EncodeFrame encodes a frame header (GEUL §2.3). When the body word count is
// below the 1023 sentinel it fits a single header word: [0001 | Class | Len].
// A body of 1023 words or more sets Len to the 1023 sentinel and appends a
// NUM32 carrying the real body word count, yielding three words in total.
func EncodeFrame(f Frame) []Word {
	if f.Len < frameLenSentinel {
		w := Word(framePrefix<<12 | uint16(f.Class)<<10 | uint16(f.Len))
		return []Word{w}
	}
	head := Word(framePrefix<<12 | uint16(f.Class)<<10 | uint16(frameLenSentinel))
	ext := EncodeNum32(f.Len)
	return []Word{head, ext[0], ext[1]}
}
