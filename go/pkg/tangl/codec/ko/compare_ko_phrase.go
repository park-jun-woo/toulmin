//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what compareKoPhrase — reverse-lookup a KO operator predicate phrase by canonical Op and negation
package ko

// compareKoPhrase finds the KO operator phrase for a canonical Compare.Op and
// negation flag by scanning the shared compareOpsKo table, returning the predicate
// phrase plus whether it carries a value and a 와/과 object particle. It is the
// reverse of parseCompareKo's suffix match, so renderCompareKo emits exactly the
// form the parser re-absorbs. ok is false when no row matches (an unknown operator).
func compareKoPhrase(op string, negated bool) (phrase string, hasValue, objParticle, ok bool) {
	for _, o := range compareOpsKo {
		if o.Op == op && o.Negated == negated {
			return o.Ko, o.HasValue, o.ObjParticle, true
		}
	}
	return "", false, false, false
}
