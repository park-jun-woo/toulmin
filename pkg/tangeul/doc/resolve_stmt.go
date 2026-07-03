//ff:func feature=tangl type=util control=selection
//ff:what ResolveStmt — resolve a StmtRef to its AST element per the README §2 path table
package doc

import (
	"errors"
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
)

// ResolveStmt returns the AST element that ref coordinates into, implementing the
// normative StmtRef path table (README contract §2) for all seven sections. It is
// tangeul-owned and exported: the stream Writer uses it, and codecs may too. Any
// out-of-range index or malformed Path is a hard error.
func ResolveStmt(d *Doc, ref StmtRef) (any, error) {
	if d == nil || d.AST == nil {
		return nil, errors.New("tangeul/doc: ResolveStmt on nil Doc or AST")
	}
	a := d.AST
	switch ref.Sec {
	case codebook.SectionSubject:
		if len(ref.Path) != 0 {
			return nil, fmt.Errorf("tangeul/doc: Subject Path must be nil, got %v", ref.Path)
		}
		return a.Subject, nil
	case codebook.SectionSee:
		if len(ref.Path) != 1 {
			return nil, fmt.Errorf("tangeul/doc: See Path must be [i], got %v", ref.Path)
		}
		if ref.Path[0] < 0 || ref.Path[0] >= len(a.Sees) {
			return nil, fmt.Errorf("tangeul/doc: See index %d out of range (%d)", ref.Path[0], len(a.Sees))
		}
		return a.Sees[ref.Path[0]], nil
	case codebook.SectionDefinitions:
		if len(ref.Path) < 1 || len(ref.Path) > 2 {
			return nil, fmt.Errorf("tangeul/doc: Definitions Path must be [i] or [i,j], got %v", ref.Path)
		}
		if ref.Path[0] < 0 || ref.Path[0] >= len(a.Defs) {
			return nil, fmt.Errorf("tangeul/doc: Definitions index %d out of range (%d)", ref.Path[0], len(a.Defs))
		}
		def := a.Defs[ref.Path[0]]
		if len(ref.Path) == 1 {
			return def, nil
		}
		if ref.Path[1] < 0 || ref.Path[1] >= len(def.Fields) {
			return nil, fmt.Errorf("tangeul/doc: Definition field index %d out of range (%d)", ref.Path[1], len(def.Fields))
		}
		return def.Fields[ref.Path[1]], nil
	case codebook.SectionRules:
		if len(ref.Path) != 1 {
			return nil, fmt.Errorf("tangeul/doc: Rules Path must be [i], got %v", ref.Path)
		}
		if ref.Path[0] < 0 || ref.Path[0] >= len(a.Rules) {
			return nil, fmt.Errorf("tangeul/doc: Rules index %d out of range (%d)", ref.Path[0], len(a.Rules))
		}
		return a.Rules[ref.Path[0]], nil
	case codebook.SectionCases:
		if len(ref.Path) != 1 && len(ref.Path) != 3 {
			return nil, fmt.Errorf("tangeul/doc: Cases Path must be [i] or [i,sk,si], got %v", ref.Path)
		}
		if ref.Path[0] < 0 || ref.Path[0] >= len(a.Cases) {
			return nil, fmt.Errorf("tangeul/doc: Cases index %d out of range (%d)", ref.Path[0], len(a.Cases))
		}
		if len(ref.Path) == 1 {
			return a.Cases[ref.Path[0]], nil
		}
		return resolveCaseSlot(&a.Cases[ref.Path[0]], ref.Path[1], ref.Path[2])
	case codebook.SectionProvides:
		if len(ref.Path) != 1 && len(ref.Path) != 3 {
			return nil, fmt.Errorf("tangeul/doc: Provides Path must be [i] or [i,sk,si], got %v", ref.Path)
		}
		if ref.Path[0] < 0 || ref.Path[0] >= len(a.Provides) {
			return nil, fmt.Errorf("tangeul/doc: Provides index %d out of range (%d)", ref.Path[0], len(a.Provides))
		}
		if len(ref.Path) == 1 {
			return a.Provides[ref.Path[0]], nil
		}
		return resolveEndpointSlot(&a.Provides[ref.Path[0]], ref.Path[1], ref.Path[2])
	case codebook.SectionInternal:
		if len(ref.Path) != 1 && len(ref.Path) != 3 {
			return nil, fmt.Errorf("tangeul/doc: Internal Path must be [i] or [i,sk,si], got %v", ref.Path)
		}
		if ref.Path[0] < 0 || ref.Path[0] >= len(a.Internals) {
			return nil, fmt.Errorf("tangeul/doc: Internal index %d out of range (%d)", ref.Path[0], len(a.Internals))
		}
		if len(ref.Path) == 1 {
			return a.Internals[ref.Path[0]], nil
		}
		return resolveInternalSlot(&a.Internals[ref.Path[0]], ref.Path[1], ref.Path[2])
	default:
		return nil, fmt.Errorf("tangeul/doc: unknown SectionCode %d", uint16(ref.Sec))
	}
}
