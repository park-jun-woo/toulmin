//ff:type feature=tangl type=model
//ff:what dec — the Reader's mutable decoding state (AST under construction + resolution maps)
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// dec is the state threaded through the Reader. ast and blocks are the Doc under
// construction; texts maps stream-assigned Text TIDs to their strings (textNext
// issues them in packet order, matching the Writer); nodeName/nodeType recover a
// node's spelling and kind for reference resolution; funcRef resolves
// materialized function references; caseByTID/epByTID/trByTID/defByTID map a
// declaring node's TID to its AST index for routing body statements. curCase is
// the open case index (-1 outside a case). checkFixups/execCaseFixups defer the
// two forward case-name references (`checking`, `run`) until every case name is
// known.
type dec struct {
	ast    *ast.Document
	blocks []doc.Block

	texts    map[uint16]string
	textNext uint16
	nodeName map[uint16]string
	nodeType map[uint16]codebook.NodeType
	funcRef  map[uint16]funcInfo

	caseByTID map[uint16]int
	epByTID   map[uint16]int
	trByTID   map[uint16]int
	defByTID  map[uint16]int

	curCase        int
	checkFixups    []nameFixup
	execCaseFixups []nameFixup
}
