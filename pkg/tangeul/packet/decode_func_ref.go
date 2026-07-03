//ff:func feature=tangl type=util control=sequence
//ff:what decodeFuncRef — decode two words into a FuncRef (Alias TID, VerbCode)
package packet

import (
	"errors"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// decodeFuncRef decodes a FuncRef from the first two words of ws (Alias TID then
// VerbCode). Label is left empty — it is restored from a separate label Text
// packet by the stream layer. It errors if fewer than two words are present.
func decodeFuncRef(ws []word.Word) (FuncRef, error) {
	if len(ws) < 2 {
		return FuncRef{}, errors.New("tangeul/packet: FuncRef requires 2 words")
	}
	return FuncRef{Alias: TID(ws[0]), Verb: uint16(ws[1])}, nil
}
