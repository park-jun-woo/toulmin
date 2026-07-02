//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestTickFuncName — tests tickFuncName for runs-seed, checks-seed, and default-seed branches
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestTickFuncName(t *testing.T) {
	tests := []struct {
		name string
		in   ast.Internal
		idx  int
		want string
	}{
		{"runs seed used", ast.Internal{Runs: []string{"caseA"}, Checks: []string{"caseB"}}, 0, "tickCaseA0"},
		{"checks seed used when no runs", ast.Internal{Checks: []string{"caseB"}}, 1, "tickCaseB1"},
		{"default tick seed when neither present", ast.Internal{}, 2, "tickTick2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tickFuncName(tt.in, tt.idx)
			if got != tt.want {
				t.Errorf("tickFuncName() = %q, want %q", got, tt.want)
			}
		})
	}
}
