//ff:func feature=tangl type=util control=sequence
//ff:what encodeFuncRef — encode a FuncRef as its two words (Alias TID, VerbCode)
package packet

import "github.com/park-jun-woo/toulmin/pkg/tangeul/word"

// encodeFuncRef encodes f as its two-word form: Alias TID then VerbCode (§4).
// Label is not emitted — it rides as a separate label Text packet at the stream
// layer (README §6).
func encodeFuncRef(f FuncRef) []word.Word {
	return []word.Word{word.Word(f.Alias), word.Word(f.Verb)}
}
