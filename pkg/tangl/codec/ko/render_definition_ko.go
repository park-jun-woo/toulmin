//ff:func feature=tangl type=parser control=sequence
//ff:what renderDefinitionKo — render a 정의 entry head (const `는 <v> 을 뜻한다[, 타입 <ref>]` or struct `는 다음을 뜻한다`)
package ko

import (
	"encoding/json"
	"strconv"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderDefinitionKo renders a tangl:정의 entry back to its head form: a struct is the
// bare "`x` 는 다음을 뜻한다" head (its 가진다 fields render as separate StmtBlocks), while
// a constant is "`x` 는 <value> 을 뜻한다" plus an optional ", 타입 <ref>" spec clause. The
// 는(은/는) and 을(을/를) particles are 받침-selected; either allomorph re-parses the same.
func renderDefinitionKo(d ast.Definition) string {
	head := "`" + d.Name + "` " + josaFor(d.Name, "은", "는")
	if d.Kind == ast.StructDef {
		return head + " 다음을 뜻한다"
	}
	if d.Kind == ast.StringListDef {
		b, _ := json.Marshal(d.Strings)
		return head + " " + string(b) + " 을 뜻한다"
	}
	if d.Kind == ast.FileSpecDef {
		return head + " `" + d.LoadType + "` 을 " + strconv.Quote(d.FilePath) + " 파일에서 읽는다"
	}
	s := head + " " + d.Value + " " + josaFor(d.Value, "을", "를") + " 뜻한다"
	if d.SpecRef != nil {
		s += ", 타입 " + renderRefKo(*d.SpecRef)
	}
	return s
}
