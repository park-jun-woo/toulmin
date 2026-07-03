//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderStmt — RenderStmt dispatches every ResolveStmt element kind to its KO renderer and propagates a ResolveStmt error
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderStmt checks every switch case of RenderStmt against a live ResolveStmt
// result — the two bare-string shapes (Subject, and a Provides/Internal run/check
// slot), each non-string AST element kind, and the propagated ResolveStmt error when
// the Doc has no AST. The fixture doc assembles one ast.Document exercising every
// StmtRef-addressable element kind: a Subject, a See, a struct Definition with one
// Field, an InlineRule, one Case carrying a Require/Node/Attack/Exec, a Provides
// Endpoint with one Run and one Check string, and an Internal trigger with one Run
// and one Check string.
func TestRenderStmt(t *testing.T) {
	a := &ast.Document{
		Subject: "환불 정책",
		Sees:    []ast.See{{Alias: "센서", Package: "robot/sensor"}},
		Defs: []ast.Definition{{
			Name: "신청인", Kind: ast.StructDef,
			Fields: []ast.Field{{Name: "연소득", Type: "Currency"}},
		}},
		Rules: []ast.InlineRule{{
			Name: "신용 충분",
			Cond: ast.Compare{Field: "신용점수", Op: "is at least", Value: "650"},
		}},
		Cases: []ast.Case{{
			Name:     "환불 가능",
			Requires: []ast.Require{{Field: "사용자"}},
			Nodes:    []ast.Node{{Name: "인증", Role: ast.GeneralRule}},
			Attacks:  []ast.Attack{{Target: "환불 가능", Attacker: "컵 없음"}},
			Execs:    []ast.Exec{{Kind: ast.RunExec, Case: "입금 가능", Node: "잔액 충분"}},
		}},
		Provides: []ast.Endpoint{{
			Name: "환불", Runs: []string{"환불 가능"}, Checks: []string{"승인 가능"},
		}},
		Internals: []ast.Internal{{
			Kind: ast.OnEvent, Event: "버튼 눌림",
			Runs: []string{"환불 가능"}, Checks: []string{"승인 가능"},
		}},
	}
	d := &doc.Doc{AST: a}
	codec := KoCodec{}

	cases := []struct {
		name string
		ref  doc.StmtRef
		want string
	}{
		{"Subject", doc.StmtRef{Sec: codebook.SectionSubject}, "이 문서는 `환불 정책` 이다"},
		{"See", doc.StmtRef{Sec: codebook.SectionSee, Path: []int{0}}, "`센서` 는 `robot/sensor` 에서 참조한다"},
		{"Definition", doc.StmtRef{Sec: codebook.SectionDefinitions, Path: []int{0}}, "`신청인` 은 다음을 뜻한다"},
		{"Field", doc.StmtRef{Sec: codebook.SectionDefinitions, Path: []int{0, 0}}, "`연소득` 을 가진다, 타입 Currency"},
		{"InlineRule", doc.StmtRef{Sec: codebook.SectionRules, Path: []int{0}}, "`신용점수` 가 650 이상이면 `신용 충분` 이다"},
		{"Case", doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0}}, "`환불 가능` 의 경우에"},
		{"Require", doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 0, 0}}, "`사용자` 는 필수다"},
		{"Node", doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 1, 0}}, "`인증` 은 일반 규칙이다"},
		{"Attack", doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 2, 0}}, "`컵 없음` 이면 `환불 가능` 무효다"},
		{"Exec", doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 3, 0}}, "`잔액 충분` 이면 `입금 가능` 진행한다"},
		{"Endpoint", doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{0}}, "`환불` 을 한다면"},
		{"EndpointRun", doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{0, 1, 0}}, "`환불 가능` 진행한다"},
		{"EndpointCheck", doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{0, 2, 0}}, "`승인 가능` 판정한다"},
		{"Internal", doc.StmtRef{Sec: codebook.SectionInternal, Path: []int{0}}, "`버튼 눌림` 발생 시"},
		{"InternalRun", doc.StmtRef{Sec: codebook.SectionInternal, Path: []int{0, 1, 0}}, "`환불 가능` 진행한다"},
		{"InternalCheck", doc.StmtRef{Sec: codebook.SectionInternal, Path: []int{0, 2, 0}}, "`승인 가능` 판정한다"},
	}
	for _, c := range cases {
		got, err := codec.RenderStmt(d, c.ref)
		if err != nil {
			t.Errorf("%s: RenderStmt: %v", c.name, err)
			continue
		}
		if got != c.want {
			t.Errorf("%s: RenderStmt = %q, want %q", c.name, got, c.want)
		}
	}

	t.Run("ResolveError", func(t *testing.T) {
		if _, err := codec.RenderStmt(&doc.Doc{}, doc.StmtRef{Sec: codebook.SectionSubject}); err == nil {
			t.Error("RenderStmt accepted a Doc with a nil AST")
		}
	})
}
