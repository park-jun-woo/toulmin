//ff:func feature=tangl type=codegen control=iteration dimension=2
//ff:what equivEvalLines — runs the cross-test document through eval.CompileDoc over the ctx table
package gen

import (
	"encoding/json"
	"fmt"
	"sort"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/eval"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// equivEvalLines compiles the cross-test document with eval.CompileDoc (host
// functions injected via equivResolver) and evaluates every rule over every
// ctx-table entry, producing the same "<ctx>\t<rule>\t<canon>" lines the
// generated runner prints — the eval side of the eval≡gen equivalence check.
func equivEvalLines(t *testing.T, doc *ast.Document, ctxJSON string) []string {
	t.Helper()
	fns, err := eval.CompileDoc(doc, equivResolver{})
	if err != nil {
		t.Fatalf("eval.CompileDoc: %v", err)
	}
	var table []map[string]any
	if err := json.Unmarshal([]byte(ctxJSON), &table); err != nil {
		t.Fatalf("ctx table: %v", err)
	}
	specs := toulmin.Specs{&equivDiscount{MinPrice: 30}}
	names := make([]string, 0, len(fns))
	for n := range fns {
		names = append(names, n)
	}
	sort.Strings(names)
	var lines []string
	for i, m := range table {
		ctx := buildEquivCtx(m)
		for _, n := range names {
			b, ev := fns[n](ctx, specs)
			lines = append(lines, fmt.Sprintf("%d\t%s\t%s", i, n, equivCanon(b, ev)))
		}
	}
	return lines
}
