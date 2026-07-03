//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderDefinitionKo — renderDefinitionKo covers const (with/without 타입) and struct head
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderDefinitionKo checks the constant forms (bare value and value-plus-타입)
// and the struct definition head, including 을/를 and 은/는 allomorph selection.
func TestRenderDefinitionKo(t *testing.T) {
	bare := renderDefinitionKo(ast.Definition{Name: "신용 기준", Kind: ast.ConstDef, Value: "650"})
	if bare != "`신용 기준` 은 650 을 뜻한다" {
		t.Errorf("bare const = %q", bare)
	}
	withType := renderDefinitionKo(ast.Definition{
		Name: "서빙 온도", Kind: ast.ConstDef, Value: "65°C",
		SpecRef: &ast.Ref{Alias: "센서", Name: "Temperature"},
	})
	if withType != "`서빙 온도` 는 65°C 를 뜻한다, 타입 `센서`.`Temperature`" {
		t.Errorf("const with 타입 = %q", withType)
	}
	strukt := renderDefinitionKo(ast.Definition{Name: "신청인", Kind: ast.StructDef})
	if strukt != "`신청인` 은 다음을 뜻한다" {
		t.Errorf("struct head = %q", strukt)
	}
}
