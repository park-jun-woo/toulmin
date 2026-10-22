//ff:func feature=tangl type=parser control=iteration dimension=2
//ff:what parseSupersedeChainItem — parse and desugar an English supersede priority chain
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseSupersedeChainItem parses "supersede in order: `A` > `B` > `C`" and
// expands every earlier-to-later pair in canonical order.
func parseSupersedeChainItem(it item, path string) ([]ast.Supersede, bool, error) {
	rest, ok := takeKeyword(it.Text, "supersede in order:")
	if !ok {
		return nil, false, nil
	}
	names, err := parseSupersedeChainNames(rest, path, it.Line)
	if err != nil {
		return nil, true, err
	}
	return expandSupersedeChain(names, it.Line), true, nil
}

func parseSupersedeChainNames(s, path string, line int) ([]string, error) {
	var names []string
	seen := map[string]bool{}
	rest := strings.TrimSpace(s)
	for {
		name, tail, ok := takeBacktick(rest)
		if !ok {
			return nil, errAt(path, line, "expected backtick-quoted name in supersede chain")
		}
		if seen[name] {
			return nil, errAt(path, line, "duplicate name %q in supersede chain", name)
		}
		seen[name] = true
		names = append(names, name)
		rest = strings.TrimSpace(tail)
		if rest == "" {
			break
		}
		if !strings.HasPrefix(rest, ">") {
			return nil, errAt(path, line, "expected '>' in supersede chain, got %q", rest)
		}
		rest = strings.TrimSpace(rest[1:])
	}
	if len(names) < 3 {
		return nil, errAt(path, line, "supersede chain requires at least 3 names")
	}
	return names, nil
}

func expandSupersedeChain(names []string, line int) []ast.Supersede {
	edges := make([]ast.Supersede, 0, len(names)*(len(names)-1)/2)
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			edges = append(edges, ast.Supersede{Superseder: names[i], Superseded: names[j], Line: line})
		}
	}
	if len(edges) > 0 {
		edges[0].Chain = append([]string(nil), names...)
	}
	return edges
}
