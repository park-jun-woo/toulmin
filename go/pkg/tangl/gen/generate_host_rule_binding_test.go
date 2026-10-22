//ff:what TestGenerateHostRuleBinding — a bare case binding inherits host function and With from Rules
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/parser"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/validate"
)

func TestGenerateHostRuleBinding(t *testing.T) {
	src := "## tangl:Subject\n- this document is `gate`\n\n" +
		"## tangl:See\n- see `ratchet` from `example/ratchet`\n\n" +
		"## tangl:Definitions\n- `max attempts` means 3 as `ratchet`.`NewMaxAttempts`\n\n" +
		"## tangl:Rules\n- `tries exhausted` uses `ratchet`.`triesExhausted` with `max attempts`\n\n" +
		"## tangl:Cases\n- in case of `gate`\n  - `tries exhausted` is a counter rule, level fail\n"
	doc, err := parser.ParseSource(src, "gate.md")
	if err != nil {
		t.Fatal(err)
	}
	if err := validate.Validate(doc); err != nil {
		t.Fatal(err)
	}
	out, err := Generate(doc)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "g.Counter(ratchet.triesExhausted).With(maxAttemptsSpec).Level(toulmin.Fail)") {
		t.Fatalf("generated binding missing:\n%s", out)
	}
}
