//ff:func feature=tangl type=parser control=sequence
//ff:what TestApplyNodeUsingKo — (alias-qualified) 판정/checking and function-ref/using paths, plus their error branches
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestApplyNodeUsingKo(t *testing.T) {
	t.Run("CheckingOk", func(t *testing.T) {
		n := &ast.Node{}
		if err := applyNodeUsingKo(n, "`판정케이스` 판정", 1, "t.md"); err != nil {
			t.Fatalf("err=%v", err)
		}
		if n.Checking == nil || n.Checking.Name != "판정케이스" || n.Using != nil {
			t.Fatalf("got %+v", n)
		}
	})

	t.Run("CheckingMissingBacktick", func(t *testing.T) {
		n := &ast.Node{}
		if err := applyNodeUsingKo(n, "그냥텍스트 판정", 1, "t.md"); err == nil {
			t.Fatal("expected error for a missing case-name backtick before '판정'")
		}
	})

	t.Run("CheckingTrailingText", func(t *testing.T) {
		n := &ast.Node{}
		if err := applyNodeUsingKo(n, "`X` 여분 판정", 1, "t.md"); err == nil {
			t.Fatal("expected error for unexpected text before '판정'")
		}
	})

	t.Run("UsingOk", func(t *testing.T) {
		n := &ast.Node{}
		if err := applyNodeUsingKo(n, "`fn`", 1, "t.md"); err != nil {
			t.Fatalf("err=%v", err)
		}
		if n.Using == nil || n.Using.Name != "fn" || n.Checking != nil {
			t.Fatalf("got %+v", n)
		}
	})

	t.Run("UsingMissingBacktick", func(t *testing.T) {
		n := &ast.Node{}
		if err := applyNodeUsingKo(n, "노백틱함수", 1, "t.md"); err == nil {
			t.Fatal("expected error for a non-backtick function reference")
		}
	})

	t.Run("UsingTrailingText", func(t *testing.T) {
		n := &ast.Node{}
		if err := applyNodeUsingKo(n, "`fn` 여분", 1, "t.md"); err == nil {
			t.Fatal("expected error for unexpected trailing text after the function reference")
		}
	})

	t.Run("AliasCheckingOk", func(t *testing.T) {
		n := &ast.Node{}
		if err := applyNodeUsingKo(n, "`대출`.`갱신 가능` 판정", 1, "t.md"); err != nil {
			t.Fatalf("err=%v", err)
		}
		if n.Checking == nil || n.Checking.Alias != "대출" || n.Checking.Name != "갱신 가능" || n.Using != nil {
			t.Fatalf("got %+v", n.Checking)
		}
	})
}
