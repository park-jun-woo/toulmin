//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestReadErrors — ReadDoc rejects corrupted documents at every layer
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/core"
)

// TestReadErrors corrupts a valid encoded document: empty input, odd length,
// truncation (DOC length mismatch), a TAG prefix outside 000, a wrong spec
// version, a HEAD without its stamps, an unclassifiable BODY record and a BODY
// not opened by SOURCE.
func TestReadErrors(t *testing.T) {
	b, err := WriteDoc(buildTransferDoc())
	if err != nil {
		t.Fatalf("WriteDoc: %v", err)
	}
	mut := func(i int, v byte) []byte {
		c := append([]byte(nil), b...)
		c[i] = v
		return c
	}
	doc := func(inner ...core.Word) []byte {
		ws := core.WrapDocument(append(inner, 0x180F))
		return core.PutWords(ws)
	}
	head, _ := emitHead("", 1)
	bad := map[string][]byte{
		"empty":         nil,
		"odd length":    b[:len(b)-1],
		"truncated":     b[:len(b)-2],
		"bad prefix":    mut(0, 0xE0),
		"spec version":  mut(10, 'X'),
		"no stamps":     doc(0x180E),
		"foreign codec": doc(append(head, 0x1066)...),
		"no source":     doc(append(head, 0x1060|1<<7, 1, 0)...),
	}
	for name, in := range bad {
		if _, err := ReadDoc(in); err == nil {
			t.Errorf("%s: expected error, got nil", name)
		}
	}
}
