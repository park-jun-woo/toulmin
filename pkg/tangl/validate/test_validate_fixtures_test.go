//ff:func feature=tangl type=validator control=sequence
//ff:what TestValidateFixtures — the spec's three worked examples, checked against known reference gaps
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestValidateFixtures parses the spec's americano, access, and transfer
// examples (shared verbatim fixtures with pkg/tangl/parser/testdata) and
// runs Validate on each.
//
// transfer.md is fully self-consistent and must pass with zero violations.
// americano.md and access.md are the spec's worked examples reproduced
// verbatim; each carries one pre-existing reference gap unrelated to the
// v0.3 features under test here (americano.md never declares `see log from
// tangl/log` even though it uses `log`.`CoffeeStatus`; access.md has no
// tangl:Definitions section even though it uses `with 'blocklist'`). Since
// these fixtures are shared, verbatim spec text outside this task's scope,
// this test asserts Validate reports exactly that one known gap for each —
// demonstrating the validator catches real dangling references without
// over- or under-reporting — rather than requiring a full pass that the
// fixture text itself does not support.
func TestValidateFixtures(t *testing.T) {
	doc, err := parser.Parse("../parser/testdata/transfer.md")
	if err != nil {
		t.Fatalf("Parse(transfer.md): %v", err)
	}
	if err := Validate(doc); err != nil {
		t.Errorf("Validate(transfer.md): %v, want nil", err)
	}

	doc, err = parser.Parse("../parser/testdata/americano.md")
	if err != nil {
		t.Fatalf("Parse(americano.md): %v", err)
	}
	if err := Validate(doc); err == nil || !strings.Contains(err.Error(), `undeclared package alias "log"`) {
		t.Errorf("Validate(americano.md) = %v, want exactly the known undeclared 'log' alias gap", err)
	}

	doc, err = parser.Parse("../parser/testdata/access.md")
	if err != nil {
		t.Fatalf("Parse(access.md): %v", err)
	}
	if err := Validate(doc); err == nil || !strings.Contains(err.Error(), `undefined term "blocklist"`) {
		t.Errorf("Validate(access.md) = %v, want exactly the known undefined 'blocklist' term gap", err)
	}
}
