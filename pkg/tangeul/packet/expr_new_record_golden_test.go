//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestNewExprRecordGolden — the three Phase003 record kinds emit their fixed word columns (header bit-layout anchor)
package packet

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestNewExprRecordGolden pins the exact word column of each Phase003 EXPR
// record — ExprCall (0x6000), ExprExists (0x2001) and ExprCompareX (0xA000) —
// so the 3-bit header bit layout (§design 1) and the operand/value body order
// (§design 2/3) are a regression anchor. Header words encode kind in bits 15-14
// plus bit 13; ExprExists/CompareX also exercise the operand sub-record and
// ValueOperand tag.
func TestNewExprRecordGolden(t *testing.T) {
	cases := []struct {
		name   string
		in     ast.Expr
		want   []word.Word
		labels []string
		n      int
	}{
		{
			name: "call-alias",
			in:   ast.Call{Fn: ast.Ref{Alias: "check", Name: "isSpam"}},
			// [ExprCall][alias "check"=0][name "isSpam"=1]
			want:   []word.Word{0x6000, 0x0000, 0x0001},
			labels: []string{"check", "isSpam"},
			n:      1,
		},
		{
			name: "exists",
			in:   ast.Exists{Source: ast.Path{Steps: []string{"coupons"}}, Binder: "c", Cond: ast.Compare{Field: "x", Op: "equals", Value: "1"}},
			// [ExprExists,count1][QuantAny][OperandPath][step1]["coupons"=0]["c"=1] + Cond kind-0 Compare
			want:   []word.Word{0x2001, 0x0321, 0x0000, 0x0001, 0x0000, 0x0001, 0x0000, 0x0303, 0x0002, 0x0000, 0x0000, 0x0001},
			labels: []string{"coupons", "c", "x"},
			n:      2,
		},
		{
			name: "comparex-field-lhs-operand-rhs",
			in:   ast.Compare{Field: "ip", Op: "equals", Rhs: ast.Path{Steps: []string{"clientIp"}}},
			// [ExprCompareX][equals]["ip"=0][lhsPresent 0][ValueOperand][OperandPath][step1]["clientIp"=1]
			want:   []word.Word{0xA000, 0x0303, 0x0000, 0x0000, 0x0002, 0x0000, 0x0001, 0x0001},
			labels: []string{"ip", "clientIp"},
			n:      1,
		},
		{
			name: "comparex-path-lhs-text-value",
			in:   ast.Compare{Op: "is at least", Value: "0.7", Lhs: ast.Path{Steps: []string{"author", "trustScore"}}},
			// [ExprCompareX][is at least][""=0][lhsPresent 1][OperandPath][step2]["author"=1]["trustScore"=2][ValueText]["0.7"=3]
			want:   []word.Word{0xA000, 0x0308, 0x0000, 0x0001, 0x0000, 0x0002, 0x0001, 0x0002, 0x0001, 0x0003},
			labels: []string{"", "author", "trustScore", "0.7"},
			n:      1,
		},
	}
	for _, c := range cases {
		var labels []string
		got, n, err := emitExpr(c.in, &labels)
		if err != nil || n != c.n {
			t.Fatalf("%s: emitExpr n=%d err=%v, want n=%d", c.name, n, err, c.n)
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Fatalf("%s: records = %v, want %v", c.name, got, c.want)
		}
		if !reflect.DeepEqual(labels, c.labels) {
			t.Fatalf("%s: labels = %q, want %q", c.name, labels, c.labels)
		}
	}
}
