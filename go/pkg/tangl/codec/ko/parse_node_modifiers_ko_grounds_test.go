//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseNodeModifiersKoGrounds — parseNodeModifiersKo 근거 modifier before 가중치, and an out-of-range 가중치
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestParseNodeModifiersKoGrounds(t *testing.T) {
	t.Run("GroundsThenQualified", func(t *testing.T) {
		n := &ast.Node{}
		rest, err := parseNodeModifiersKo(n, "`body` 그리고 `mx` 를 근거로 하는, 가중치 0.5 의 `검사`.`F` 쓰는", 1, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if len(n.Grounds) != 2 || n.Grounds[0] != "body" || n.Grounds[1] != "mx" || n.Qualified == nil || *n.Qualified != 0.5 {
			t.Fatalf("grounds=%v qualified=%v", n.Grounds, n.Qualified)
		}
		if rest != "`검사`.`F` 쓰는" {
			t.Fatalf("rest=%q", rest)
		}
	})

	t.Run("QualifiedOutOfRange", func(t *testing.T) {
		if _, err := parseNodeModifiersKo(&ast.Node{}, "가중치 1.5 의", 1, "t.md"); err == nil {
			t.Fatal("expected an out-of-range 가중치 error")
		}
	})
}
