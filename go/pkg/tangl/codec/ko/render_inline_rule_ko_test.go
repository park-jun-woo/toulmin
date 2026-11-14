//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderInlineRuleKo — renderInlineRuleKo covers the one-line and 다음이면 nested forms
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderInlineRuleKo checks a single-compare rule (one line) and a nested rule
// whose Logic/Either tree renders as a "다음이면" head plus a 다음 중/그리고/또는 condition
// list at relative indentation.
func TestRenderInlineRuleKo(t *testing.T) {
	single := ast.InlineRule{
		Name: "신용 충분",
		Cond: ast.Compare{Field: "신용점수", Op: "is at least", Value: "650"},
	}
	got, err := renderInlineRuleKo(single)
	if err != nil {
		t.Fatalf("single: %v", err)
	}
	if got != "`신용점수` 가 650 이상이면 `신용 충분` 이다" {
		t.Errorf("single rule = %q", got)
	}

	nested := ast.InlineRule{
		Name: "우대 대상",
		Cond: ast.Logic{Op: "or", Terms: []ast.Expr{
			ast.Either{Terms: []ast.Expr{
				ast.Logic{Op: "and", Terms: []ast.Expr{
					ast.Compare{Field: "등급", Op: "equals", Value: `"VIP"`},
					ast.Compare{Field: "연체 이력", Op: "is empty"},
				}},
			}},
			ast.Compare{Field: "직원", Op: "equals", Value: "true"},
		}},
	}
	got, err = renderInlineRuleKo(nested)
	if err != nil {
		t.Fatalf("nested: %v", err)
	}
	want := "다음이면 `우대 대상` 이다\n" +
		"  - 다음 중\n" +
		"    - `등급` 이 \"VIP\" 와 같으면\n" +
		"    - 그리고 `연체 이력` 이 비어 있으면\n" +
		"  - 또는 `직원` 이면"
	if got != want {
		t.Errorf("nested rule =\n%q\nwant\n%q", got, want)
	}

	predErr := ast.InlineRule{
		Name: "우대 대상",
		Cond: ast.Compare{Field: "x", Op: "bogus", Value: "1"},
	}
	if _, err := renderInlineRuleKo(predErr); err == nil {
		t.Error("renderInlineRuleKo accepted a Compare with an unknown operator")
	}

	listErr := ast.InlineRule{
		Name: "우대 대상",
		Cond: ast.Logic{Op: "and", Terms: []ast.Expr{
			ast.Compare{Field: "직원", Op: "equals", Value: "true"},
			ast.Compare{Field: "등급", Op: "equals", Value: `"VIP"`},
			ast.Compare{Field: "연체 이력", Op: "is empty"},
		}},
	}
	if _, err := renderInlineRuleKo(listErr); err == nil {
		t.Error("renderInlineRuleKo accepted a non-binary Logic condition list")
	}
}
