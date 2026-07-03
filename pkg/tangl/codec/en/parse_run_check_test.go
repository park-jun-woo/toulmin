//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseRunCheck — classify run/check/other body lines and reject malformed backtick names
package en

import "testing"

// TestParseRunCheck covers all five outcomes: a `run` line (slotKind 1), a
// `check` line (slotKind 2), a line matching neither keyword (isBody false, no
// error), a run/check line missing its opening backtick, and one with an
// unterminated backtick-quoted name.
func TestParseRunCheck(t *testing.T) {
	slotKind, name, isBody, err := parseRunCheck("run `a`")
	if slotKind != 1 || name != "a" || !isBody || err != nil {
		t.Errorf("run: got (%d, %q, %v, %v), want (1, \"a\", true, nil)", slotKind, name, isBody, err)
	}

	slotKind, name, isBody, err = parseRunCheck("check `b`")
	if slotKind != 2 || name != "b" || !isBody || err != nil {
		t.Errorf("check: got (%d, %q, %v, %v), want (2, \"b\", true, nil)", slotKind, name, isBody, err)
	}

	slotKind, name, isBody, err = parseRunCheck("`q` is required")
	if slotKind != 0 || name != "" || isBody || err != nil {
		t.Errorf("other: got (%d, %q, %v, %v), want (0, \"\", false, nil)", slotKind, name, isBody, err)
	}

	if _, _, isBody, err = parseRunCheck("run oops"); !isBody || err == nil {
		t.Errorf("missing backtick: isBody=%v, err=%v, want (true, non-nil)", isBody, err)
	}

	if _, _, isBody, err = parseRunCheck("run `oops"); !isBody || err == nil {
		t.Errorf("unterminated backtick: isBody=%v, err=%v, want (true, non-nil)", isBody, err)
	}
}
