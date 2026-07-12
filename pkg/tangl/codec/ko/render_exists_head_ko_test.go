//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderExistsHeadKo — nested exists head surface (다음인 것이 있으면) and source-operand error
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderExistsHeadKo(t *testing.T) {
	t.Run("Head", func(t *testing.T) {
		x := ast.Exists{Source: ast.Path{Steps: []string{"쿠폰들"}}, Binder: "항목"}
		got, err := renderExistsHeadKo(x)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if want := "`쿠폰들` 중 `항목` 이 다음인 것이 있으면"; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("OpenVowelBinder", func(t *testing.T) {
		x := ast.Exists{Source: ast.Path{Steps: []string{"쿠폰들"}}, Binder: "c"}
		got, err := renderExistsHeadKo(x)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if want := "`쿠폰들` 중 `c` 가 다음인 것이 있으면"; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("SourceError", func(t *testing.T) {
		x := ast.Exists{Source: &ast.Path{Steps: []string{"쿠폰들"}}, Binder: "c"}
		if _, err := renderExistsHeadKo(x); err == nil {
			t.Fatal("expected error for an unsupported source operand")
		}
	})
}
