//ff:func feature=tangl type=util control=sequence
//ff:what EqualIgnoringLine — compare two ASTs for equality while ignoring Line source positions
package tangeul

import (
	"reflect"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// EqualIgnoringLine reports whether two ASTs are structurally identical when the
// `Line` source-position field on every node is disregarded. A .geul stream
// does not store authored line numbers (Phase 006: the reader leaves Line = 0),
// so the retargeting contract only requires equality up to Line. The comparison
// treats a nil and an empty slice as equal (a representational, not semantic,
// difference). It is exported because tangl-en and tangl-ko reuse it to check
// their encoders against the direct text-parse path.
func EqualIgnoringLine(a, b *ast.Document) bool {
	return equalIgnoringLine(reflect.ValueOf(a), reflect.ValueOf(b))
}
