//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRoundtripReencode — L3 mode-2 self-consistency: re-encoding the normalized rendering yields the same AST
package en

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// TestRoundtripReencode pins mode-2 normalization self-consistency (Phase 004 L3):
// for each fixture it encodes the source (d1), renders the EN canonical form with
// doc.Render, re-encodes that rendering (d2), and asserts d2's AST equals d1's up
// to Line. A normalized rendering fed back through the codec must reproduce the
// same graph — the fixed point that makes mode-2 output re-executable. This runs
// wholly in memory (no stream Save), so americano is included alongside the other
// seven fixtures.
func TestRoundtripReencode(t *testing.T) {
	td := func(name string) string {
		return filepath.Join("..", "..", "parser", "testdata", name)
	}
	ex := func(name string) string {
		return filepath.Join("..", "..", "..", "..", "..", "test", name)
	}
	paths := []string{
		td("americano.md"), td("access.md"), td("transfer.md"),
		td("rules_expr.md"),
		ex("library-renewal.md"), ex("newsletter-tick-warn.md"),
		ex("newsletter-tick-fixed.md"), ex("order-saga.md"),
		ex("factory-line-cascade.md"),
	}
	for _, p := range paths {
		src, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("read %s: %v", p, err)
		}
		d1, err := ENCodec{}.Encode(src, p)
		if err != nil {
			t.Errorf("%s: Encode: %v", p, err)
			continue
		}
		rendered, err := doc.Render(d1, ENCodec{})
		if err != nil {
			t.Errorf("%s: Render: %v", p, err)
			continue
		}
		d2, err := ENCodec{}.Encode(rendered, p)
		if err != nil {
			t.Errorf("%s: re-encode of rendered output failed: %v\n---\n%s", p, err, rendered)
			continue
		}
		if !tangeul.EqualIgnoringLine(d1.AST, d2.AST) {
			t.Errorf("%s: re-encoded AST differs from original\n--- rendered ---\n%s", p, rendered)
		}
	}
}
