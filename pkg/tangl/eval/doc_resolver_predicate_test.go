//ff:func feature=tangl type=engine control=sequence
//ff:what TestDocResolverPredicate — tests sibling-rule-first predicate resolution vs the host
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestDocResolverPredicate(t *testing.T) {
	host := fakeResolver{predicates: map[string]func(toulmin.Context, toulmin.Specs) (bool, any){
		"nlp\x00spam": func(toulmin.Context, toulmin.Specs) (bool, any) { return true, "host" },
	}}
	d := docResolver{
		host:    host,
		ruleSet: map[string]bool{"r": true},
		compiled: map[string]evalFunc{
			"r": func(sc *scope) (bool, any) { return true, "rule" },
		},
	}

	fn, ok := d.Predicate("", "r")
	if !ok {
		t.Fatal("bare rule name should resolve to the sibling rule")
	}
	if b, w := fn(testCtx(), nil); !b || w != "rule" {
		t.Fatalf("sibling rule closure: got %v,%v", b, w)
	}
	fn, ok = d.Predicate("nlp", "spam")
	if !ok {
		t.Fatal("aliased call should delegate to the host")
	}
	if b, w := fn(testCtx(), nil); !b || w != "host" {
		t.Fatalf("host predicate: got %v,%v", b, w)
	}
	if _, ok := d.Predicate("", "unknown"); ok {
		t.Fatal("bare non-rule name should miss via the host")
	}
}
