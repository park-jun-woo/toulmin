//ff:func feature=tangl type=util control=sequence
//ff:what fakeLoader — test helper returning a loadAST closure over a fixed path->Document map
package registry

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// fakeLoader returns a loadAST closure over a fixed path->Document map,
// stamping Path onto each returned document (as tangeul.Load does) and
// counting how many times each path is parsed. It lets loadRegistry's
// recursion be exercised without real .tangeul fixtures.
func fakeLoader(docs map[string]*ast.Document, counts map[string]int) func(string) (*ast.Document, error) {
	return func(path string) (*ast.Document, error) {
		d, ok := docs[path]
		if !ok {
			return nil, fmt.Errorf("fakeLoader: no document at %q", path)
		}
		counts[path]++
		d.Path = path
		return d, nil
	}
}
