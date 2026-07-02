//ff:func feature=tangl type=parser control=sequence
//ff:what applyCaseChild — dispatch and append one case-body statement
package parser

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// applyCaseChild tries each case-body statement form in turn (required
// field, node registration, don't attack edge, do/undo/run exec edge) and
// appends the match onto c.
func applyCaseChild(c *ast.Case, child item, path string) error {
	if req, ok, err := parseRequireItem(child, path); err != nil {
		return err
	} else if ok {
		c.Requires = append(c.Requires, req)
		return nil
	}
	if node, ok, err := parseNodeItem(child, path); err != nil {
		return err
	} else if ok {
		c.Nodes = append(c.Nodes, node)
		return nil
	}
	if atk, ok, err := parseAttackItem(child, path); err != nil {
		return err
	} else if ok {
		c.Attacks = append(c.Attacks, atk)
		return nil
	}
	if exec, ok, err := parseExecItem(child, path); err != nil {
		return err
	} else if ok {
		c.Execs = append(c.Execs, exec)
		return nil
	}
	return errAt(path, child.Line, "unrecognized case statement: %q", child.Text)
}
