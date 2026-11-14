//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestApplyCaseChildKo — each slot kind (require/node/attack/exec) appends and reports its index
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestApplyCaseChildKo(t *testing.T) {
	c := &ast.Case{Name: "c"}
	steps := []struct {
		text string
		kind int
	}{
		{"`사용자` 는 필수다", 0},
		{"`인증` 은 일반 규칙이다", 1},
		{"`컵 없음` 이면 `주문 접수` 무효다", 2},
		{"`주문 접수` 면 `물 붓기 가능` 진행한다", 3},
	}
	for _, s := range steps {
		sk, err := applyCaseChildKo(c, scan.Item{Text: s.text, No: 1}, "t.md")
		if err != nil {
			t.Fatalf("%q: %v", s.text, err)
		}
		if sk != s.kind {
			t.Errorf("%q: slotKind=%d, want %d", s.text, sk, s.kind)
		}
	}
	if len(c.Requires) != 1 || len(c.Nodes) != 1 || len(c.Attacks) != 1 || len(c.Execs) != 1 {
		t.Fatalf("appended counts: req=%d node=%d atk=%d exec=%d", len(c.Requires), len(c.Nodes), len(c.Attacks), len(c.Execs))
	}
	if _, err := applyCaseChildKo(c, scan.Item{Text: "알 수 없는 문장", No: 2}, "t.md"); err == nil {
		t.Fatal("expected error for an unrecognized case statement")
	}

	// Each callee (parseRequireKo/parseNodeKo/parseAttackKo/parseExecKo) can itself
	// return ok=true with a non-nil error (as opposed to ok=false with no error);
	// applyCaseChildKo must propagate that error immediately instead of falling
	// through to the next case-body form.
	calleeErrors := []struct {
		name string
		text string
	}{
		{"RequireErr", "`f` 는 필수다 나머지"}, // parseRequireKo: tail after 필수다 isn't '타입'
		{"NodeErr", "`X` 규칙이다"},         // parseNodeKo: no role word before 규칙이다
		{"AttackErr", "이상한 문장 무효다"},     // parseAttackKo: no backtick-quoted attacker
		{"ExecErr", "이상한 문장 진행한다"},      // parseExecKo/parseRunExecKo: no backtick-quoted node
	}
	for _, ce := range calleeErrors {
		if _, err := applyCaseChildKo(c, scan.Item{Text: ce.text, No: 3}, "t.md"); err == nil {
			t.Errorf("%s: %q: expected a propagated callee error", ce.name, ce.text)
		}
	}
}
