//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what decodeLabels — read an EXPR label table (count word + UTF8-TEXT Values) and the words consumed
package packet

import (
	"errors"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
)

// decodeLabels reads the EXPR label table at the head of ws — a count word, then
// that many UTF8-TEXT Value records — and returns the strings with the number
// of words consumed. It is the inverse of encodeLabels.
func decodeLabels(ws []core.Word) ([]string, int, error) {
	if len(ws) < 1 {
		return nil, 0, errors.New("tangeul/packet: EXPR label table requires a count word")
	}
	pos := 1
	labels := make([]string, 0, int(ws[0]))
	for i := 0; i < int(ws[0]); i++ {
		v, err := toulmin.TakeValue(ws, &pos)
		if err != nil {
			return nil, 0, err
		}
		b, err := core.DecodeBytesValue(v, core.CodecUTF8TextEven)
		if err != nil {
			return nil, 0, err
		}
		labels = append(labels, string(b))
	}
	return labels, pos, nil
}
