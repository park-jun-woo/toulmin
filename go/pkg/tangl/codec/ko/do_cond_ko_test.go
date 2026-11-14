//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestDoCondKo — doCondKo renders a 받침-selected 면/이면 with no gate, delegates to renderCertaintyKo otherwise, and propagates its error
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestDoCondKo checks the two doCondKo branches: e.Certainty == nil returns the
// 받침-selected bare 면/이면 for e.Node (covering both allomorphs), and a non-nil
// Certainty delegates to renderCertaintyKo, including its error path for an unknown
// operator.
func TestDoCondKo(t *testing.T) {
	cases := []struct {
		name    string
		e       ast.Exec
		want    string
		wantErr bool
	}{
		{"NoGateBatchim", ast.Exec{Node: "팔"}, "이면", false},
		{"NoGateOpen", ast.Exec{Node: "휴가"}, "면", false},
		{"GateOk", ast.Exec{Node: "팔", Certainty: &ast.Certainty{Op: "at least", Percent: 80}}, "확신 80% 이상이면", false},
		{"GateErr", ast.Exec{Node: "팔", Certainty: &ast.Certainty{Op: "bogus", Percent: 1}}, "", true},
	}
	for _, c := range cases {
		got, err := doCondKo(c.e)
		gotErr := err != nil
		if gotErr != c.wantErr {
			t.Errorf("%s: doCondKo(%+v) err = %v, wantErr %v", c.name, c.e, err, c.wantErr)
			continue
		}
		if c.wantErr {
			continue
		}
		if got != c.want {
			t.Errorf("%s: doCondKo(%+v) = %q, want %q", c.name, c.e, got, c.want)
		}
	}
}
