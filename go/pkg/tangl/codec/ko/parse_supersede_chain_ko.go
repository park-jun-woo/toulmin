//ff:func feature=tangl type=parser control=iteration dimension=2
//ff:what parseSupersedeChainKo — parse and desugar a Korean supersede priority chain
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func parseSupersedeChainKo(it scan.Item, path string) ([]ast.Supersede, bool, error) {
	head, ok := takeSuffixKo(it.Text, "순서로 배제한다")
	if !ok {
		return nil, false, nil
	}
	var names []string
	seen := map[string]bool{}
	rest := strings.TrimSpace(head)
	for {
		name, tail, ok := takeBacktick(rest)
		if !ok {
			return nil, true, errAtKo(path, it.No, "expected backtick-quoted name in supersede chain")
		}
		if seen[name] {
			return nil, true, errAtKo(path, it.No, "duplicate name %q in supersede chain", name)
		}
		seen[name] = true
		names = append(names, name)
		rest = strings.TrimSpace(tail)
		if rest == "" {
			break
		}
		if !strings.HasPrefix(rest, ">") {
			return nil, true, errAtKo(path, it.No, "expected '>' in supersede chain, got %q", rest)
		}
		rest = strings.TrimSpace(rest[1:])
	}
	if len(names) < 3 {
		return nil, true, errAtKo(path, it.No, "supersede chain requires at least 3 names")
	}
	edges := make([]ast.Supersede, 0, len(names)*(len(names)-1)/2)
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			edges = append(edges, ast.Supersede{Superseder: names[i], Superseded: names[j], Line: it.No})
		}
	}
	if len(edges) > 0 {
		edges[0].Chain = append([]string(nil), names...)
	}
	return edges, true, nil
}
