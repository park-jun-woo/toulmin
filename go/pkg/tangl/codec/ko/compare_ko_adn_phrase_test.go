//ff:func feature=tangl type=parser control=sequence
//ff:what TestCompareKoAdnPhrase — every table row round-trips 관형형 phrase ⇄ (Op, negated) both directions
package ko

import (
	"strings"
	"testing"
)

func TestCompareKoAdnPhrase(t *testing.T) {
	t.Run("ReverseLookupMatchesTable", func(t *testing.T) {
		for _, o := range compareOpsKo {
			phrase, hasValue, objParticle, ok := compareKoAdnPhrase(o.Op, o.Negated)
			if !ok {
				t.Fatalf("no adnominal phrase for Op=%q negated=%v", o.Op, o.Negated)
			}
			if phrase != o.KoAdn || hasValue != o.HasValue || objParticle != o.ObjParticle {
				t.Fatalf("Op=%q negated=%v: got (%q,%v,%v), want (%q,%v,%v)",
					o.Op, o.Negated, phrase, hasValue, objParticle, o.KoAdn, o.HasValue, o.ObjParticle)
			}
		}
	})

	t.Run("AdnominalEndsInAttributive", func(t *testing.T) {
		// Every 관형형 form ends in the attributive ㄴ/은/는/인 (not the 서술형 ~면).
		for _, o := range compareOpsKo {
			if strings.HasSuffix(o.KoAdn, "면") {
				t.Errorf("KoAdn %q still ends in the 서술형 ~면", o.KoAdn)
			}
		}
	})

	t.Run("Unknown", func(t *testing.T) {
		if _, _, _, ok := compareKoAdnPhrase("no such op", false); ok {
			t.Fatal("expected ok=false for an unknown operator")
		}
	})
}
