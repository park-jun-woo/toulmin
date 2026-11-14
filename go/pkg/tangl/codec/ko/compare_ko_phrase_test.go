//ff:func feature=tangl type=parser control=sequence
//ff:what TestCompareKoPhrase — compareKoPhrase reverse-looks-up every table row and reports ok=false for an unmatched op/negation
package ko

import "testing"

// TestCompareKoPhrase checks that every row of compareOpsKo is found by its exact
// Op/Negated pair (ok=true, exact phrase/hasValue/objParticle echoed back) and that an
// unmatched pair — an unknown Op, and a known Op with the wrong Negated value — returns
// the zero value with ok=false.
func TestCompareKoPhrase(t *testing.T) {
	t.Run("MatchesEveryRow", func(t *testing.T) {
		for _, o := range compareOpsKo {
			phrase, hasValue, objParticle, ok := compareKoPhrase(o.Op, o.Negated)
			if !ok {
				t.Fatalf("compareKoPhrase(%q, %v) ok = false, want true", o.Op, o.Negated)
			}
			if phrase != o.Ko || hasValue != o.HasValue || objParticle != o.ObjParticle {
				t.Errorf("compareKoPhrase(%q, %v) = (%q, %v, %v), want (%q, %v, %v)",
					o.Op, o.Negated, phrase, hasValue, objParticle, o.Ko, o.HasValue, o.ObjParticle)
			}
		}
	})

	t.Run("UnknownOp", func(t *testing.T) {
		phrase, hasValue, objParticle, ok := compareKoPhrase("is unknown", false)
		if ok || phrase != "" || hasValue || objParticle {
			t.Errorf("compareKoPhrase(unknown) = (%q, %v, %v, %v), want (\"\", false, false, false)", phrase, hasValue, objParticle, ok)
		}
	})

	t.Run("WrongNegation", func(t *testing.T) {
		// "is empty" only has a Negated=false row; requesting Negated=true must miss.
		phrase, hasValue, objParticle, ok := compareKoPhrase("is empty", true)
		if ok || phrase != "" || hasValue || objParticle {
			t.Errorf("compareKoPhrase(\"is empty\", true) = (%q, %v, %v, %v), want (\"\", false, false, false)", phrase, hasValue, objParticle, ok)
		}
	})
}
