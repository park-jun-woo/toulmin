//ff:func feature=tangl type=util control=sequence
//ff:what TestReadErrors — corrupt streams (truncation, odd length, bad prefix/version, no source) all error
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestReadErrors checks that malformed input fails loudly rather than silently:
// an empty stream, an odd byte length, a truncated final packet, a frame prefix
// violation, an unknown GEUL major version, and a header-only stream with no
// Source payload are each a hard error.
func TestReadErrors(t *testing.T) {
	b, err := WriteDoc(buildTransferDoc())
	if err != nil {
		t.Fatalf("WriteDoc: %v", err)
	}

	corrupt := func(name string, in []byte) {
		if _, err := ReadDoc(in); err == nil {
			t.Errorf("%s: expected error, got nil", name)
		}
	}
	corrupt("empty", nil)
	corrupt("odd length", b[:len(b)-1])
	corrupt("truncated packet", b[:len(b)-2])

	badPrefix := append([]byte(nil), b...)
	badPrefix[0] = 0x0C // frame prefix nibble 0000 instead of 0001 (Free)
	corrupt("bad frame prefix", badPrefix)

	badVersion := append([]byte(nil), b...)
	badVersion[4] = 0x09 // GEUL major byte in the version stamp
	corrupt("unknown major", badVersion)

	corrupt("header only, no source", word.PutWords(emitHeader(codebook.LocaleUnknown)))
}
