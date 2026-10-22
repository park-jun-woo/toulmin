//ff:func feature=tangl type=util control=sequence
//ff:what TestHostRuleRoundTrip — Rules host rules and in-case nodes keep using/with/grounds through .geul
package tangeul

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/en"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/parser"
)

// TestHostRuleRoundTrip parses a document whose tangl:Rules declares host rules
// with `with` specs and `grounded on` facts, saves it as .geul, loads it back and
// requires the same AST — the F_FUNC/F_WITH/F_GROUNDS path of the TOULMIN codec.
func TestHostRuleRoundTrip(t *testing.T) {
	src := "## tangl:Subject\n- this document is `mail`\n\n## tangl:See\n- see `check` from `x/check`\n\n" +
		"## tangl:Definitions\n- `limit` means 5\n\n" +
		"## tangl:Rules\n- `ip blocked` uses `check`.`IsBlocked` with `limit`\n" +
		"- `source lacks email` uses `check`.`SourceLacks` grounded on `source-body` and `mx`\n\n" +
		"## tangl:Cases\n- in case of `can send`\n  - `ip blocked` is a counter rule\n  - `source lacks email` is a counter rule\n"
	dir := t.TempDir()
	md := filepath.Join(dir, "mail.md")
	if err := os.WriteFile(md, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	a, err := parser.Parse(md)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(a.Rules) != 2 || len(a.Rules[1].Grounds) != 2 || len(a.Rules[0].With) != 1 {
		t.Fatalf("fixture did not parse as host rules: %+v", a.Rules)
	}
	enc, err := en.ENCodec{}.Encode([]byte(src), md)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	out := filepath.Join(dir, "mail.geul")
	if err := Save(enc, out); err != nil {
		t.Fatalf("save: %v", err)
	}
	d, err := Load(out)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	d.AST.Path = a.Path
	if !EqualIgnoringLine(a, d.AST) {
		t.Fatalf("host-rule AST drifted:\n got %+v\nwant %+v", d.AST.Rules, a.Rules)
	}
}
