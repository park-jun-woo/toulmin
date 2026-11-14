//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseForceAttackItem — explicit-force attack edge: success (incl. force 1), not-this-form, and every error branch
package parser

import (
	"strings"
	"testing"
)

func TestParseForceAttackItem(t *testing.T) {
	a, ok, err := parseForceAttackItem(item{Text: "`block ip` attacks `authenticate` with force 0.5", Line: 3}, "t.md")
	if err != nil || !ok || a.Attacker != "block ip" || a.Target != "authenticate" || a.Force == nil || *a.Force != 0.5 || a.Line != 3 {
		t.Fatalf("got (%+v, %v, %v)", a, ok, err)
	}
	a, ok, err = parseForceAttackItem(item{Text: "`b` attacks `a` with force 1.0", Line: 4}, "t.md")
	if err != nil || !ok || a.Force == nil || *a.Force != 1 {
		t.Fatalf("explicit 1.0 must stay explicit: got (%+v, %v, %v)", a, ok, err)
	}
	for _, text := range []string{"do `f` when `n`", "`b` is a general rule", "`b` supersedes `a`"} {
		if _, ok, err := parseForceAttackItem(item{Text: text, Line: 1}, "t.md"); ok || err != nil {
			t.Errorf("parseForceAttackItem(%q) = ok %v err %v, want not-this-form", text, ok, err)
		}
	}
	errs := []struct {
		text, want string
	}{
		{"`b` attacks a with force 0.5", "expected backtick-quoted target"},
		{"`b` attacks `a` force 0.5", "expected 'with force'"},
		{"`b` attacks `a` with force 2", "out of range"},
		{"`b` attacks `a` with force 0.5 now", "unexpected trailing text"},
	}
	for _, tc := range errs {
		_, ok, err := parseForceAttackItem(item{Text: tc.text, Line: 9}, "t.md")
		if !ok || err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("parseForceAttackItem(%q) = ok %v err %v, want %q", tc.text, ok, err, tc.want)
		}
	}
}
