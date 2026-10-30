//ff:func feature=tangl type=engine control=sequence
//ff:what TestDefError — the formatted message is prefixed with path, line and definition name
package materialize

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestDefError(t *testing.T) {
	err := defError(&ast.Document{Path: "p.md"}, ast.Definition{Name: "d", Line: 2}, "loader %q is not registered", "X")
	if err.Error() != `p.md:2: definition "d": loader "X" is not registered` {
		t.Fatalf("defError = %v", err)
	}
}
