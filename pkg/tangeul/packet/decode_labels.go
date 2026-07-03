//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what decodeLabels — read an EXPR label table from ws, returning the labels and words consumed
package packet

import (
	"errors"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// decodeLabels reads the EXPR label table from the head of ws (a count word,
// then per label a NUM32 byte length and its padded UTF-8 words) and returns the
// decoded strings with the number of words consumed. It is the inverse of
// encodeLabels.
func decodeLabels(ws []word.Word) ([]string, int, error) {
	if len(ws) < 1 {
		return nil, 0, errors.New("tangeul/packet: EXPR label table requires a count word")
	}
	count := int(ws[0])
	pos := 1
	labels := make([]string, count)
	for i := 0; i < count; i++ {
		if len(ws) < pos+2 {
			return nil, 0, errors.New("tangeul/packet: EXPR label length truncated")
		}
		byteLen, err := word.DecodeNum32(ws[pos : pos+2])
		if err != nil {
			return nil, 0, err
		}
		pos += 2
		dataWords := int((byteLen + 1) / 2)
		if len(ws) < pos+dataWords {
			return nil, 0, errors.New("tangeul/packet: EXPR label bytes truncated")
		}
		labels[i] = string(word.PutWords(ws[pos : pos+dataWords])[:byteLen])
		pos += dataWords
	}
	return labels, pos, nil
}
