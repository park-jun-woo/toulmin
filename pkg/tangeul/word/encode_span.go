//ff:func feature=tangl type=util control=sequence
//ff:what EncodeSpan — encode a Span as 3 words (TID, offset, length)
package word

// EncodeSpan encodes s as three words in the fixed order TID, offset, length
// (GEUL §5.5).
func EncodeSpan(s Span) [spanWords]Word {
	return [spanWords]Word{
		Word(s.TID),
		Word(s.Off),
		Word(s.Len),
	}
}
