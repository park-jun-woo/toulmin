//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCompareKo — 8 operators × positive/negative, allomorph invariance, longest-match, error
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseCompareKo(t *testing.T) {
	t.Run("PositiveOperators", func(t *testing.T) {
		cases := []struct {
			in    string
			field string
			op    string
			value string
		}{
			{"`f` 가 비어 있으면", "f", "is empty", ""},
			{"`f` 가 비어 있지 않으면", "f", "is not empty", ""},
			{"`등급` 이 \"VIP\" 와 같으면", "등급", "equals", "VIP"},
			{"`f` 가 목록 중 하나면", "f", "is in", "목록"},
			{"`등급` 이 \"VIP\", \"GOLD\" 중 하나면", "등급", "is in", "VIP, GOLD"},
			{"`나이` 가 20 보다 크면", "나이", "is greater than", "20"},
			{"`나이` 가 20 보다 작으면", "나이", "is less than", "20"},
			{"`신용점수` 가 650 이상이면", "신용점수", "is at least", "650"},
			{"`점수` 가 90 이하면", "점수", "is at most", "90"},
		}
		for _, c := range cases {
			expr, err := parseCompareKo(c.in, 1, "t.md")
			if err != nil {
				t.Fatalf("%q: unexpected error: %v", c.in, err)
			}
			cmp, ok := expr.(ast.Compare)
			if !ok {
				t.Fatalf("%q: got %T, want Compare", c.in, expr)
			}
			if cmp.Field != c.field || cmp.Op != c.op || cmp.Value != c.value {
				t.Fatalf("%q: got %+v", c.in, cmp)
			}
		}
	})

	t.Run("NegativeOperators", func(t *testing.T) {
		cases := []struct {
			in string
			op string
		}{
			{"`등급` 이 \"VIP\" 와 같지 않으면", "equals"},
			{"`f` 가 `목록` 중 하나가 아니면", "is in"},
			{"`나이` 가 20 보다 크지 않으면", "is greater than"},
			{"`나이` 가 20 보다 작지 않으면", "is less than"},
			{"`신용점수` 가 650 이상이 아니면", "is at least"},
			{"`점수` 가 90 이하가 아니면", "is at most"},
		}
		for _, c := range cases {
			expr, err := parseCompareKo(c.in, 1, "t.md")
			if err != nil {
				t.Fatalf("%q: unexpected error: %v", c.in, err)
			}
			not, ok := expr.(ast.Not)
			if !ok {
				t.Fatalf("%q: got %T, want Not", c.in, expr)
			}
			cmp, ok := not.Term.(ast.Compare)
			if !ok || cmp.Op != c.op {
				t.Fatalf("%q: inner = %+v", c.in, not.Term)
			}
		}
	})

	t.Run("AllomorphInvariance", func(t *testing.T) {
		a, err1 := parseCompareKo("`등급` 이 \"VIP\" 와 같으면", 1, "t.md")
		b, err2 := parseCompareKo("`등급` 가 \"VIP\" 과 같으면", 1, "t.md")
		if err1 != nil || err2 != nil {
			t.Fatalf("errs: %v %v", err1, err2)
		}
		if a != b {
			t.Fatalf("allomorph forms differ: %+v vs %+v", a, b)
		}
	})

	t.Run("Unrecognized", func(t *testing.T) {
		if _, err := parseCompareKo("`f` 가 3 이랑 비슷하면", 1, "t.md"); err == nil {
			t.Fatal("expected error for an unrecognized operator")
		}
	})

	t.Run("MissingBacktickField", func(t *testing.T) {
		_, err := parseCompareKo("필드 가 비어 있으면", 1, "t.md")
		if err == nil {
			t.Fatal("expected error for a field not wrapped in backticks")
		}
		const want = `t.md:1: expected backtick-quoted field in "필드 가 비어 있으면"`
		if got := err.Error(); got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}
