//ff:func feature=cli type=command control=sequence
//ff:what loadExecutable — load a .geul for execution, rejecting text .md surfaces with guidance
package tanglcli

import (
	"fmt"
	"path/filepath"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// loadExecutable is the single entry point the executing commands (check, ast,
// effects, gen) use to read their input. TANGEUL is the only executable surface
// (README §5), so a text `.md` path is rejected up front with a message pointing
// the user at `tangl encode`; anything else is handed to tangeul.Load, whose own
// extension check enforces `.geul`.
func loadExecutable(path string) (*doc.Doc, error) {
	if filepath.Ext(path) == ".md" {
		return nil, fmt.Errorf("text surface is not executable; run 'tangl encode %s' first", path)
	}
	return tangeul.Load(path)
}
