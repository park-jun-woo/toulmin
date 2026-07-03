//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what PutWords — serialize words to bytes, big-endian (high byte first)
package word

// PutWords serializes ws to bytes in big-endian order (high byte first per
// word), producing 2 bytes per word (GEUL §2.1).
func PutWords(ws []Word) []byte {
	out := make([]byte, 0, len(ws)*2)
	for _, w := range ws {
		out = append(out, byte(w>>8), byte(w&0xFF))
	}
	return out
}
