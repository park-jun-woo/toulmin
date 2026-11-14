//ff:func feature=tangl type=engine control=sequence
//ff:what TestDocResolverPredicateNoHost — without a host Resolver only sibling rules resolve
package eval

import "testing"

func TestDocResolverPredicateNoHost(t *testing.T) {
	d := docResolver{ruleSet: map[string]bool{"r": true}, compiled: map[string]evalFunc{
		"r": func(sc *scope) (bool, any) { return false, "sibling" },
	}}
	if fn, ok := d.Predicate("nlp", "spam"); ok || fn != nil {
		t.Fatalf("aliased call with no host = (%v, %v), want (nil, false)", fn != nil, ok)
	}
	if fn, ok := d.Predicate("", "unknown"); ok || fn != nil {
		t.Fatal("bare non-rule name with no host must miss")
	}
	fn, ok := d.Predicate("", "r")
	if !ok {
		t.Fatal("sibling rule must resolve even without a host")
	}
	if got, ev := fn(testCtx(), nil); got || ev != "sibling" {
		t.Fatalf("sibling closure = %v,%v", got, ev)
	}
}
