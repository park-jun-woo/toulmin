//ff:func feature=tangl type=util control=sequence
//ff:what newDec — build a fresh Reader state with empty maps and no open case
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// newDec builds an initialized Reader state: an empty AST, the resolution maps,
// a Text-TID counter starting at 1 (matching the Writer), and curCase -1 (outside
// any case).
func newDec() *dec {
	return &dec{
		ast:       &ast.Document{},
		texts:     map[uint16]string{},
		textNext:  1,
		nodeName:  map[uint16]string{},
		nodeType:  map[uint16]codebook.NodeType{},
		funcRef:   map[uint16]funcInfo{},
		caseByTID: map[uint16]int{},
		epByTID:   map[uint16]int{},
		trByTID:   map[uint16]int{},
		defByTID:  map[uint16]int{},
		curCase:   -1,
	}
}
