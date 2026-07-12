//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what compareKoAdnPhrase — reverse-lookup a KO 관형형 operator phrase by canonical Op and negation
package ko

// compareKoAdnPhrase finds the KO 관형형 (attributive) operator phrase for a
// canonical Compare.Op and negation flag by scanning the shared compareOpsKo table,
// returning the attributive phrase plus whether it carries a value and a 와/과 object
// particle. It is the 관형형 mirror of compareKoPhrase (which returns the 서술형 Ko
// column), so renderCompareAdnKo emits exactly the form parseCompareAdnKo re-absorbs.
// ok is false when no row matches (an unknown operator).
func compareKoAdnPhrase(op string, negated bool) (phrase string, hasValue, objParticle, ok bool) {
	for _, o := range compareOpsKo {
		if o.Op == op && o.Negated == negated {
			return o.KoAdn, o.HasValue, o.ObjParticle, true
		}
	}
	return "", false, false, false
}
