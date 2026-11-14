//ff:func feature=tangl type=util control=sequence
//ff:what loadRegistry — DFS-loads every checking-alias-reachable document, rejecting file-load cycles
package registry

import (
	"fmt"
	"path/filepath"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// loadRegistry is the loader-injectable core of Load: it DFS-walks the
// checking-alias reference graph starting at rootPath, using loadAST to read
// each document's AST. A canonical absolute path already on the current DFS
// stack is a file-load cycle (rejected regardless of whether a case-level
// checking cycle would manifest); a path already in reg.Docs but off the
// stack is a completed diamond node and is reused, not re-parsed.
func loadRegistry(rootPath string, loadAST func(string) (*ast.Document, error)) (*Registry, error) {
	abs, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}
	reg := &Registry{Docs: make(map[string]*ast.Document)}
	stack := make(map[string]bool)
	var visit func(path string) error
	visit = func(path string) error {
		if stack[path] {
			return fmt.Errorf("registry: document load cycle detected at %q", path)
		}
		if _, ok := reg.Docs[path]; ok {
			return nil
		}
		d, err := loadAST(path)
		if err != nil {
			return err
		}
		reg.Docs[path] = d
		stack[path] = true
		for _, c := range d.Cases {
			for _, n := range c.Nodes {
				if n.Checking == nil || n.Checking.Alias == "" {
					continue
				}
				see, ok := findSee(d, n.Checking.Alias)
				if !ok {
					return fmt.Errorf("registry: %s: checking alias %q has no matching See declaration", path, n.Checking.Alias)
				}
				child, err := resolveCheckingAliasPath(d.Path, see.Package)
				if err != nil {
					return err
				}
				if err := visit(child); err != nil {
					return err
				}
			}
		}
		delete(stack, path)
		return nil
	}
	if err := visit(abs); err != nil {
		return nil, err
	}
	return reg, nil
}
