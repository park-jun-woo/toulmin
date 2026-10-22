//ff:func feature=tangl type=parser control=sequence
//ff:what applyCaseChildKo — dispatch one case-body item, append it, and report its slot kind
package ko

import (
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// applyCaseChildKo tries each case-body statement form in turn — required field
// (필수다), node registration (규칙이다), don't attack edge (무효다), supersede edge
// (배제한다), then do/undo/run exec edge (실행한다/되돌린다/진행한다) — appends the match
// onto c, and returns the StmtRef slot kind (0 Require · 1 Node · 2 Attack · 3 Exec ·
// 4 Supersede) so the encoder can assign the body coordinate. It mirrors
// pkg/tangl/parser.applyCaseChild.
func applyCaseChildKo(c *ast.Case, it scan.Item, path string) (int, error) {
	if req, ok, err := parseRequireKo(it, path); err != nil {
		return 0, err
	} else if ok {
		c.Requires = append(c.Requires, req)
		return 0, nil
	}
	if node, ok, err := parseNodeKo(it, path); err != nil {
		return 0, err
	} else if ok {
		c.Nodes = append(c.Nodes, node)
		return 1, nil
	}
	if atk, ok, err := parseAttackKo(it, path); err != nil {
		return 0, err
	} else if ok {
		c.Attacks = append(c.Attacks, atk)
		return 2, nil
	}
	if edges, ok, err := parseSupersedeChainKo(it, path); err != nil {
		return 0, err
	} else if ok {
		c.Supersedes = append(c.Supersedes, edges...)
		return 4, nil
	}
	if sup, ok, err := parseSupersedeKo(it, path); err != nil {
		return 0, err
	} else if ok {
		c.Supersedes = append(c.Supersedes, sup)
		return 4, nil
	}
	if exec, ok, err := parseExecKo(it, path); err != nil {
		return 0, err
	} else if ok {
		c.Execs = append(c.Execs, exec)
		return 3, nil
	}
	return 0, errAtKo(path, it.No, "unrecognized case statement: %q", it.Text)
}
