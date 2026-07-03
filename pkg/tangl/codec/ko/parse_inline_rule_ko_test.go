//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseInlineRuleKo — one-line form, nested 다음이면 form, EN↔KO AST equivalence, missing 이다/name
package ko

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

func TestParseInlineRuleKo(t *testing.T) {
	t.Run("OneLine", func(t *testing.T) {
		rule, err := parseInlineRuleKo(scan.Item{Text: "`신용점수` 가 650 이상이면 `신용 충분` 이다", No: 1}, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rule.Name != "신용 충분" {
			t.Fatalf("name = %q", rule.Name)
		}
		want := ast.Compare{Field: "신용점수", Op: "is at least", Value: "650"}
		if rule.Cond != want {
			t.Fatalf("cond = %#v, want %#v", rule.Cond, want)
		}
	})

	t.Run("MissingIda", func(t *testing.T) {
		if _, err := parseInlineRuleKo(scan.Item{Text: "`신용점수` 가 650 이상이면 `신용 충분`", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error for missing 이다")
		}
	})

	t.Run("MissingName", func(t *testing.T) {
		if _, err := parseInlineRuleKo(scan.Item{Text: "650 이상이면 이다", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error for missing backtick rule name")
		}
	})

	t.Run("OneLineCondError", func(t *testing.T) {
		if _, err := parseInlineRuleKo(scan.Item{Text: "바나나 `이름` 이다", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error propagated from a malformed one-line condition")
		}
	})

	t.Run("NestedCondError", func(t *testing.T) {
		if _, err := parseInlineRuleKo(scan.Item{Text: "다음이면 `이름` 이다", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error propagated from an empty nested condition list")
		}
	})

	// NestedMatchesEN parses the grammar §규칙 nested example and asserts the
	// condition tree is byte-for-byte the AST the EN parser produces from the
	// mirrored EN source.
	t.Run("NestedMatchesEN", func(t *testing.T) {
		koItem := scan.Item{
			Text: "다음이면 `우대 대상` 이다",
			No:   1,
			Children: []scan.Item{
				{
					Text: "다음 중",
					No:   2,
					Children: []scan.Item{
						{Text: "`등급` 이 \"VIP\" 와 같으면", No: 3},
						{Text: "그리고 `연체 이력` 이 비어 있으면", No: 4},
					},
				},
				{Text: "또는 `직원` 이면", No: 5},
			},
		}
		rule, err := parseInlineRuleKo(koItem, "ko.md")
		if err != nil {
			t.Fatalf("KO parse error: %v", err)
		}

		enSrc := "## tangl:Subject\n- this document is `t`\n\n" +
			"## tangl:Rules\n" +
			"1. `우대 대상` when\n" +
			"    - either\n" +
			"      - `등급` equals \"VIP\"\n" +
			"      - and `연체 이력` is empty\n" +
			"    - or `직원` equals true\n\n" +
			"## tangl:Cases\n- in case of `x`\n  - `a` is required\n"
		doc, err := parser.ParseSource(enSrc, "en.md")
		if err != nil {
			t.Fatalf("EN parse error: %v", err)
		}
		if len(doc.Rules) != 1 {
			t.Fatalf("EN Rules = %+v, want 1", doc.Rules)
		}
		if !reflect.DeepEqual(rule.Cond, doc.Rules[0].Cond) {
			t.Fatalf("KO cond != EN cond\n KO: %#v\n EN: %#v", rule.Cond, doc.Rules[0].Cond)
		}
	})
}
