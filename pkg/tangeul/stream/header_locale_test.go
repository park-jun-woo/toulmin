//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestHeaderLocaleRoundtrip — the authoring locale stamped by WriteDoc is recovered by ReadDoc, and a pre-locale header degrades to unknown
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestHeaderLocaleRoundtrip serializes the transfer fixture with an authoring
// locale set on the Doc and asserts ReadDoc recovers it from the STREAM_HEADER.
// It also drives readHeaderLocale directly over a 5-word (pre-locale) body to
// confirm the backward-compatible degrade to LocaleUnknown.
func TestHeaderLocaleRoundtrip(t *testing.T) {
	for _, locale := range []string{"en", "ko", ""} {
		d := buildTransferDoc()
		d.Locale = locale
		b, err := WriteDoc(d)
		if err != nil {
			t.Fatalf("WriteDoc(%q): %v", locale, err)
		}
		got, err := ReadDoc(b)
		if err != nil {
			t.Fatalf("ReadDoc(%q): %v", locale, err)
		}
		if got.Locale != locale {
			t.Errorf("recovered Locale = %q, want %q", got.Locale, locale)
		}
	}

	// A pre-locale header carries a 5-word body (no locale word); readHeaderLocale
	// must degrade to LocaleUnknown rather than read past the body.
	body := []word.Word{word.Word(codebook.ControlStreamHeader), 0, 0, 0, 0}
	ws := append(word.EncodeFrame(word.Frame{Class: word.ClassControl, Len: uint32(len(body))}), body...)
	code, err := readHeaderLocale(rawPacket{Class: word.ClassControl, Words: ws})
	if err != nil {
		t.Fatalf("readHeaderLocale(short body): %v", err)
	}
	if code != codebook.LocaleUnknown {
		t.Errorf("readHeaderLocale(short body) = %d, want LocaleUnknown", code)
	}

	// A rawPacket whose Words cannot even be framed (empty) makes word.DecodeFrame
	// fail; readHeaderLocale must propagate that error and report LocaleUnknown.
	code, err = readHeaderLocale(rawPacket{Class: word.ClassControl, Words: nil})
	if err == nil {
		t.Fatal("readHeaderLocale(unframeable words): expected an error, got nil")
	}
	if code != codebook.LocaleUnknown {
		t.Errorf("readHeaderLocale(unframeable words) code = %d, want LocaleUnknown", code)
	}
}
