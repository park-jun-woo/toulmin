//ff:func feature=tangl type=codegen control=sequence
//ff:what withArg — resolves a "with" clause term to its .With() Go expression
package gen

// withArg resolves a "with <term>" clause to the Go expression passed to
// .With(): the term's derived Spec var if Definitions declared one with a
// SpecRef, its plain const identifier if Definitions declared it without
// one, or — if the term isn't in Definitions at all — the term itself as
// a bare identifier (assumed to be a package-level Spec value defined by
// hand or via a See-imported package; codegen cannot verify this).
func withArg(gc *genContext, term string) string {
	if info, ok := gc.Defs[term]; ok {
		if info.Spec != "" {
			return info.Spec
		}
		if info.Const != "" {
			return info.Const
		}
	}
	return goIdent(term)
}
