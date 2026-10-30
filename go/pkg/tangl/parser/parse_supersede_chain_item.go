//ff:func feature=tangl type=parser control=sequence
//ff:what parseSupersedeChainItem — parse and desugar an English supersede priority chain
package parser

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

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
