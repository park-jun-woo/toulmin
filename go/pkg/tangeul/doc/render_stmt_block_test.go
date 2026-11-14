//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestRenderStmtBlock — single-line and multi-line statements at indent level 0 and 1
package doc

import "testing"

// TestRenderStmtBlock verifies renderStmtBlock's list-marker/indent skeleton: a
// single-line statement takes only the loop's i==0 branch (list marker, no
// continuation indent), a multi-line statement also takes the else branch for
// its continuation lines, and a non-zero indentLevel prefixes every line with
// the extra base indent.
func TestRenderStmtBlock(t *testing.T) {
	tests := []struct {
		name        string
		stmt        string
		indentLevel int
		want        string
	}{
		{
			name:        "single line, level 0",
			stmt:        "`a` is 5",
			indentLevel: 0,
			want:        "- `a` is 5\n",
		},
		{
			name:        "single line, level 1",
			stmt:        "`a` is 5",
			indentLevel: 1,
			want:        "  - `a` is 5\n",
		},
		{
			name:        "multi line, level 0",
			stmt:        "`r` when\n`a` is 5\nand `b`",
			indentLevel: 0,
			want:        "- `r` when\n  `a` is 5\n  and `b`\n",
		},
		{
			name:        "multi line, level 1",
			stmt:        "`r` when\n`a` is 5",
			indentLevel: 1,
			want:        "  - `r` when\n    `a` is 5\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderStmtBlock(tt.stmt, tt.indentLevel)
			if string(got) != tt.want {
				t.Errorf("renderStmtBlock(%q, %d) = %q, want %q", tt.stmt, tt.indentLevel, got, tt.want)
			}
		})
	}
}
