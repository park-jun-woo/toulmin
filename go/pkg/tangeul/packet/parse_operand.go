//ff:func feature=tangl type=util control=selection
//ff:what parseOperand — decode one ast.Operand (Path/SpecRef/Call) from EXPR records, resolving labels
package packet

import (
	"errors"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseOperand decodes one operand sub-record from recs at *pos, advancing *pos
// past it and resolving label indices through labels (§design 2). It is the
// exact inverse of emitOperand: an OperandPath/OperandSpec reads a step count
// and that many label indices; an OperandCall reads an alias and a name label
// index. An out-of-range index or an unknown OperandKind is a hard error.
func parseOperand(recs []core.Word, pos *int, labels []string) (ast.Operand, error) {
	if len(recs) < *pos+1 {
		return nil, errors.New("tangeul/packet: EXPR operand kind missing")
	}
	kind := OperandKind(recs[*pos])
	*pos++
	switch kind {
	case OperandPath:
		steps, err := parseOperandSteps(recs, pos, labels)
		if err != nil {
			return nil, err
		}
		return ast.Path{Steps: steps}, nil
	case OperandSpec:
		steps, err := parseOperandSteps(recs, pos, labels)
		if err != nil {
			return nil, err
		}
		return ast.SpecRef{Steps: steps}, nil
	case OperandCall:
		if len(recs) < *pos+2 {
			return nil, errors.New("tangeul/packet: EXPR operand Call truncated")
		}
		alias := int(recs[*pos])
		name := int(recs[*pos+1])
		*pos += 2
		if alias >= len(labels) || name >= len(labels) {
			return nil, errors.New("tangeul/packet: EXPR operand Call label out of range")
		}
		return ast.Call{Fn: ast.Ref{Alias: labels[alias], Name: labels[name]}}, nil
	default:
		return nil, errors.New("tangeul/packet: unknown OperandKind")
	}
}
