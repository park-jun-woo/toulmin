//ff:func feature=tangl type=util control=sequence
//ff:what newDec — an initialized Reader state for a document with the given field widths
package stream

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// newDec builds an initialized Reader state: an empty AST, the document's field
// widths, empty Edge Id and resolution maps, and curCase -1 (outside any case).
func newDec(w toulmin.Wire) *dec {
	return &dec{
		ast:       &ast.Document{},
		wire:      w,
		seen:      map[toulmin.TID]bool{},
		texts:     map[toulmin.TID]string{},
		nodeName:  map[toulmin.TID]string{},
		nodeType:  map[toulmin.TID]toulmin.NodeType{},
		funcRef:   map[toulmin.TID]funcInfo{},
		caseByTID: map[toulmin.TID]int{},
		epByTID:   map[toulmin.TID]int{},
		trByTID:   map[toulmin.TID]int{},
		defByTID:  map[toulmin.TID]int{},
		curCase:   -1,
	}
}
