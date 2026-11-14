//ff:func feature=cli type=command control=sequence
//ff:what saveTangeul — serialize a Doc to a .geul path for CLI tests, failing the test on error
package tanglcli

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// saveTangeul serializes d to path via tangeul.Save (path must end in .geul)
// and fails the test on error. It is the CLI tests' bridge from a programmatic Doc
// fixture to an on-disk .geul the commands can load.
func saveTangeul(t *testing.T, path string, d *doc.Doc) {
	t.Helper()
	if err := tangeul.Save(d, path); err != nil {
		t.Fatalf("Save(%s): %v", path, err)
	}
}
