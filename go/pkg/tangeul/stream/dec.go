//ff:type feature=tangl type=model
//ff:what dec — the Reader's state: AST under construction, wire widths, Edge Id registry and resolution maps
package stream

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// dec is the state threaded through the Reader. ast and blocks are the Doc under
// construction; wire is the document's declared field widths; seen records every
// Edge Id read so a duplicate is rejected; texts maps TEXT Edge Ids to their
// strings; nodeName/nodeType recover a node's spelling and kind for reference
// resolution; funcRef resolves materialized function references;
// caseByTID/epByTID/trByTID/defByTID map a declaring node's Edge Id to its AST
// index for routing body statements. curCase is the open case index (-1 outside
// a case). checkFixups/execCaseFixups defer the two forward case-name
// references (`checking`, `run`) until every case name is known.
type dec struct {
	ast    *ast.Document
	blocks []doc.Block
	wire   toulmin.Wire
	seen   map[toulmin.TID]bool

	texts    map[toulmin.TID]string
	nodeName map[toulmin.TID]string
	nodeType map[toulmin.TID]toulmin.NodeType
	funcRef  map[toulmin.TID]funcInfo

	caseByTID map[toulmin.TID]int
	epByTID   map[toulmin.TID]int
	trByTID   map[toulmin.TID]int
	defByTID  map[toulmin.TID]int

	curCase        int
	checkFixups    []nameFixup
	execCaseFixups []nameFixup
}
