//ff:what TestQuotedLiteralIntegration — quoted comparison literals evaluate as their unquoted values
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestQuotedLiteralIntegration(t *testing.T) {
	const src = "## tangl:Subject\n- this document is `quoted`\n\n" +
		"## tangl:Rules\n" +
		"- `version matches` when `submission`.`version` equals \"v0.3.1\"\n" +
		"- `main exists` when any of `files` as `f` satisfies `f`.`name` equals \"main.go\"\n\n" +
		"- `grade listed` when `grade` is in \"VIP\", \"GOLD\"\n\n" +
		"## tangl:Cases\n" +
		"- in case of `gate`\n" +
		"  - `version matches` is a general rule\n"

	doc, err := parser.ParseSource(src, "quoted.md")
	if err != nil {
		t.Fatalf("ParseSource: %v", err)
	}
	fns, err := CompileDoc(doc, nil)
	if err != nil {
		t.Fatalf("CompileDoc: %v", err)
	}
	ctx := toulmin.NewContext()
	ctx.Set("submission", map[string]any{"version": "v0.3.1"})
	ctx.Set("files", []any{map[string]any{"name": "README.md"}, map[string]any{"name": "main.go"}})
	ctx.Set("grade", "VIP")

	for _, name := range []string{"version matches", "main exists", "grade listed"} {
		got, evidence := fns[name](ctx, nil)
		if !got {
			t.Errorf("%q = false, want true; evidence=%+v", name, evidence)
		}
	}

	ctx.Set("submission", map[string]any{"version": "v0.3.0"})
	got, evidence := fns["version matches"](ctx, nil)
	if got {
		t.Fatal("mismatched version = true, want false")
	}
	fact, ok := evidence.(toulmin.Fact)
	if !ok {
		t.Fatalf("evidence = %T, want toulmin.Fact", evidence)
	}
	if fact.Expected != "equals v0.3.1" {
		t.Fatalf("Fact.Expected = %q, want quotes removed", fact.Expected)
	}

	ctx.Set("grade", "SILVER")
	got, evidence = fns["grade listed"](ctx, nil)
	if got {
		t.Fatal("unlisted grade = true, want false")
	}
	fact, ok = evidence.(toulmin.Fact)
	if !ok || fact.Expected != "is in VIP, GOLD" {
		t.Fatalf("is-in evidence = %#v, want normalized Expected", evidence)
	}
}
