//ff:func feature=cli type=command control=selection
//ff:what runGraph — dispatches graph command by file extension
package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// runGraph dispatches by file extension: .yaml/.yml → YAML flow, .go → Go analysis.
func runGraph(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".yaml", ".yml":
		return runGraphYAML(cmd, filePath)
	case ".go":
		return runGraphGo(cmd, filePath)
	default:
		return fmt.Errorf("unsupported file extension %q: use .yaml, .yml, or .go", ext)
	}
}
