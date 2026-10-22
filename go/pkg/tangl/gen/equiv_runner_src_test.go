//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what equivRunnerSrc — builds the cross-test runner main.go source
package gen

import (
	"strconv"
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// equivRunnerHead is the fixed part of the cross-test runner: byte-identical
// copies of the shared host functions, ctx builder, spec/author types, and
// canonical formatter defined as equiv* helpers in this package.
const equivRunnerHead = `package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	toulmin "github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

type equivAuthor struct {
	TrustScore float64
}

type equivDiscount struct {
	MinPrice float64
}

func (s *equivDiscount) SpecName() string { return "discount" }

func (s *equivDiscount) Validate() error { return nil }

func nlpScore(ctx toulmin.Context, specs toulmin.Specs) (any, error) {
	if v, ok := ctx.Get("nlpErr"); ok && v == true {
		return nil, errors.New("nlp unavailable")
	}
	v, _ := ctx.Get("nlpScore")
	return v, nil
}

func isSpam(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	v, _ := ctx.Get("spamFlag")
	ev, _ := ctx.Get("spamEvidence")
	return v == true, ev
}

func buildEquivCtx(m map[string]any) toulmin.Context {
	ctx := toulmin.NewContext()
	for k, v := range m {
		if k == "authorTrustScore" {
			ctx.Set("author", &equivAuthor{TrustScore: v.(float64)})
			continue
		}
		ctx.Set(k, v)
	}
	return ctx
}

func equivCanon(b bool, ev any) string {
	if f, ok := ev.(toulmin.Fact); ok {
		return fmt.Sprintf("%t|fact|%v|%v|%T|%v", b, f.Where, f.Expected, f.Actual, f.Actual)
	}
	return fmt.Sprintf("%t|ev|%T|%v", b, ev, ev)
}

func main() {
	var table []map[string]any
	if err := json.Unmarshal([]byte(equivCtxTable), &table); err != nil {
		panic(err)
	}
	specs := toulmin.Specs{&equivDiscount{MinPrice: 30}}
	names := make([]string, 0, len(equivRules))
	for n := range equivRules {
		names = append(names, n)
	}
	sort.Strings(names)
	for i, m := range table {
		ctx := buildEquivCtx(m)
		for _, n := range names {
			b, ev := equivRules[n](ctx, specs)
			fmt.Printf("%d\t%s\t%s\n", i, n, equivCanon(b, ev))
		}
	}
}
`

// equivRunnerSrc builds the runner main.go: the fixed head plus the ctx
// table (the shared JSON, quoted) and the name-to-generated-function rule
// map derived from the document's Rules via goIdent — the same lookup the
// generated Cases graph uses.
func equivRunnerSrc(rules []ast.InlineRule, ctxJSON string) string {
	var b strings.Builder
	b.WriteString(equivRunnerHead)
	b.WriteString("\nconst equivCtxTable = " + strconv.Quote(ctxJSON) + "\n")
	b.WriteString("\nvar equivRules = map[string]func(toulmin.Context, toulmin.Specs) (bool, any){\n")
	for _, r := range rules {
		b.WriteString("\t" + strconv.Quote(r.Name) + ": " + goIdent(r.Name) + ",\n")
	}
	b.WriteString("}\n")
	return b.String()
}
