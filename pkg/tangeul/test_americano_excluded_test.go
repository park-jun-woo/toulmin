//ff:func feature=tangl type=util control=sequence
//ff:what TestAmericanoExcluded — pin the reason americano has no step-1 golden (undeclared alias)
package tangeul

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestAmericanoExcluded pins the documented step-1 exclusion: americano.md does a
// `do` on the alias `log` without a matching `see`, and the stream Writer rejects
// references to undeclared aliases, so its Save fails. americano therefore has no
// golden and stays out of the parity set. The failure is the stream Writer's, not
// the tiling's: assembleAmericanoDoc lays a valid, fully tiled Doc, and Save fails
// only when it reaches the `do log.CoffeeStatus` exec.
func TestAmericanoExcluded(t *testing.T) {
	a, err := parser.Parse(filepath.Join("..", "tangl", "parser", "testdata", "americano.md"))
	if err != nil {
		t.Fatalf("parse americano: %v", err)
	}
	d := assembleAmericanoDoc(a)
	err = Save(d, filepath.Join(t.TempDir(), "americano.tangeul"))
	if err == nil {
		t.Fatal("americano Save unexpectedly succeeded; the undeclared-alias exclusion no longer holds")
	}
	if !strings.Contains(err.Error(), `undeclared alias "log"`) {
		t.Fatalf("americano Save failed for an unexpected reason: %v", err)
	}
}
