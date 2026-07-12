//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCondTermKoExpr — dispatch of the new 판정식 branches (exists / 참이면 call / path bool) + dot-chain error
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseCondTermKoExpr(t *testing.T) {
	t.Run("ExistsOneLineDispatch", func(t *testing.T) {
		expr, err := parseCondTermKo("`쿠폰들` 중 `c` 가 `c` 의 `최소가격` 이 `기본가` 이하인 것이 있으면", nil, 1, "t.md")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if _, ok := expr.(ast.Exists); !ok {
			t.Fatalf("got %#v, want Exists", expr)
		}
	})

	t.Run("CallPredDispatch", func(t *testing.T) {
		expr, err := parseCondTermKo("`검사`.`isSpam` 가 참이면", nil, 1, "t.md")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if _, ok := expr.(ast.Call); !ok {
			t.Fatalf("got %#v, want Call", expr)
		}
	})

	t.Run("BarePathBooleanDispatch", func(t *testing.T) {
		expr, err := parseCondTermKo("`작성자` 의 `활성` 이면", nil, 1, "t.md")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		cmp, ok := expr.(ast.Compare)
		if !ok {
			t.Fatalf("got %#v, want Compare", expr)
		}
		if _, ok := cmp.Lhs.(ast.Path); !ok || cmp.Op != "equals" || cmp.Value != "true" {
			t.Fatalf("got %#v", cmp)
		}
	})

	t.Run("PathBareNegationDispatch", func(t *testing.T) {
		expr, err := parseCondTermKo("`작성자` 의 `활성` 이 아니고", nil, 1, "t.md")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		not, ok := expr.(ast.Not)
		if !ok {
			t.Fatalf("got %#v, want Not", expr)
		}
		if _, ok := not.Term.(ast.Compare).Lhs.(ast.Path); !ok {
			t.Fatalf("inner = %#v", not.Term)
		}
	})

	t.Run("UnmarkedDotChainLhsErrors", func(t *testing.T) {
		if _, err := parseCondTermKo("`a`.`b` 가 100 이상이면", nil, 1, "t.md"); err == nil {
			t.Fatal("expected a clear error for an unmarked dot-chain left side")
		}
	})
}
