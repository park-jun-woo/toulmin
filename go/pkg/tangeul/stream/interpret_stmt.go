//ff:func feature=tangl type=util control=sequence
//ff:what interpretStmt — decode one statement group into an AST element and its Block
package stream

import (
	"errors"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// interpretStmt decodes one collected statement (its prefix payloads and
// same-span Node/Edge body) into an AST element and returns the covering Block.
// It decodes the body, splits Func nodes from the primary and label edges from
// relation edges, records the prefix TEXT strings (recovering the Rules EXPR by
// position), registers Func nodes, then routes to the node- or edge-shaped
// builder.
func interpretStmt(d *dec, prefix, body []rawPacket) (doc.Block, error) {
	nodes, edges, span, err := decodeBody(d, body)
	if err != nil {
		return doc.Block{}, err
	}
	funcNodes, primaries := splitNodes(nodes)
	labelEdges, rel := splitEdges(edges)

	cond, err := assignPrefix(d, prefix)
	if err != nil {
		return doc.Block{}, err
	}
	labels, err := buildLabelMap(d, labelEdges)
	if err != nil {
		return doc.Block{}, err
	}
	registerFuncNodes(d, funcNodes, labels)

	dspan, err := fromWireSpan(span)
	if err != nil {
		return doc.Block{}, err
	}
	if len(primaries) == 1 {
		ref, err := buildFromNode(d, primaries[0], labels, rel, cond)
		if err != nil {
			return doc.Block{}, err
		}
		return doc.Block{Kind: doc.StmtBlock, Span: dspan, Ref: ref}, nil
	}
	if len(primaries) == 0 {
		ref, err := buildFromEdge(d, rel)
		if err != nil {
			return doc.Block{}, err
		}
		return doc.Block{Kind: doc.StmtBlock, Span: dspan, Ref: ref}, nil
	}
	return doc.Block{}, errors.New("tangeul/stream: statement has more than one primary node")
}
