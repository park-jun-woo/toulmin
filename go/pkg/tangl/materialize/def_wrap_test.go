//ff:func feature=tangl type=engine control=sequence
//ff:what TestDefWrap — the wrapped error carries the definition location and stays unwrappable
package materialize

import (
	"errors"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestDefWrap(t *testing.T) {
	cause := errors.New("cause")
	err := defWrap(&ast.Document{Path: "p.md"}, ast.Definition{Name: "d", Line: 7}, cause)
	if err.Error() != `p.md:7: definition "d": cause` || !errors.Is(err, cause) {
		t.Fatalf("defWrap = %v", err)
	}
}
