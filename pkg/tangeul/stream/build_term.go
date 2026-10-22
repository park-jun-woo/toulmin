//ff:func feature=tangl type=util control=sequence
//ff:what buildTerm — append a Definitions term (DefKind from role, value and spec ref from edges)
package stream

import (
	"encoding/json"
	"fmt"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// buildTerm appends the ast.Definition a Term node declares: its DefKind from the
// role bits, a const's literal value from the means edge, and an `as` spec
// reference from the sameAs edge (resolved through its Func node).
func buildTerm(d *dec, tid uint16, name string, role codebook.RoleCode, rel []packet.EdgePacket) (doc.StmtRef, error) {
	def := ast.Definition{Name: name, Kind: ast.DefKind(role)}
	if e, ok := findRel(rel, codebook.PropMeans); ok {
		value := d.texts[uint16(e.DstTID)]
		switch def.Kind {
		case ast.StringListDef:
			if err := json.Unmarshal([]byte(value), &def.Strings); err != nil {
				return doc.StmtRef{}, fmt.Errorf("decode string-list definition %q: %w", name, err)
			}
		case ast.FileSpecDef:
			var payload struct {
				LoadType string `json:"loadType"`
				FilePath string `json:"filePath"`
			}
			if err := json.Unmarshal([]byte(value), &payload); err != nil {
				return doc.StmtRef{}, fmt.Errorf("decode file definition %q: %w", name, err)
			}
			def.LoadType, def.FilePath = payload.LoadType, payload.FilePath
		default:
			def.Value = value
		}
	}
	sp, err := termSpecRef(d, rel)
	if err != nil {
		return doc.StmtRef{}, err
	}
	def.SpecRef = sp
	idx := len(d.ast.Defs)
	d.ast.Defs = append(d.ast.Defs, def)
	d.defByTID[tid] = idx
	return doc.StmtRef{Sec: codebook.SectionDefinitions, Path: []int{idx}}, nil
}
