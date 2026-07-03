//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what renderStmtBlock — shift a (possibly multi-line) statement to its section indent with a list marker
package doc

import (
	"bytes"
	"strings"
)

// renderStmtBlock places a codec statement at its section indent as a markdown
// list item: the first line gets the `- ` marker, continuation lines (a Rules
// subtree) are aligned past the marker (README contract §3). The statement's own
// relative indentation is preserved; indentLevel positions the whole block.
func renderStmtBlock(stmt string, indentLevel int) []byte {
	base := strings.Repeat("  ", indentLevel)
	var b bytes.Buffer
	for i, line := range strings.Split(stmt, "\n") {
		b.WriteString(base)
		if i == 0 {
			b.WriteString("- ")
		} else {
			b.WriteString("  ")
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.Bytes()
}
