//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseDefinitionKo — const, const+타입 spec ref, struct head with fields, missing name, missing 뜻한다
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestParseDefinitionKo(t *testing.T) {
	t.Run("ConstWithSpecRef", func(t *testing.T) {
		def, err := parseDefinitionKo(scan.Item{Text: "`신용 기준` 은 650 을 뜻한다, 타입 `credit`.`Threshold`", No: 1}, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if def.Name != "신용 기준" || def.Kind != ast.ConstDef || def.Value != "650" {
			t.Fatalf("got %+v", def)
		}
		if def.SpecRef == nil || def.SpecRef.Alias != "credit" || def.SpecRef.Name != "Threshold" {
			t.Fatalf("specRef = %+v", def.SpecRef)
		}
	})

	t.Run("ConstNoTail", func(t *testing.T) {
		def, err := parseDefinitionKo(scan.Item{Text: "`문턱` 은 10 을 뜻한다", No: 1}, "t.md")
		if err != nil || def.Kind != ast.ConstDef || def.Value != "10" || def.SpecRef != nil {
			t.Fatalf("got %+v err=%v", def, err)
		}
	})

	t.Run("StructHead", func(t *testing.T) {
		it := scan.Item{
			Text: "`신청인` 은 다음을 뜻한다",
			No:   1,
			Children: []scan.Item{
				{Text: "`연소득` 을 가진다, 타입 Currency", No: 2},
				{Text: "`신용점수` 를 가진다, 타입 Integer", No: 3},
			},
		}
		def, err := parseDefinitionKo(it, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if def.Kind != ast.StructDef || len(def.Fields) != 2 {
			t.Fatalf("got %+v", def)
		}
		if def.Fields[0].Name != "연소득" || def.Fields[0].Type != "Currency" {
			t.Fatalf("field0 = %+v", def.Fields[0])
		}
		if def.Fields[1].Name != "신용점수" || def.Fields[1].Type != "Integer" {
			t.Fatalf("field1 = %+v", def.Fields[1])
		}
	})

	t.Run("MissingName", func(t *testing.T) {
		if _, err := parseDefinitionKo(scan.Item{Text: "신용 기준 은 650 을 뜻한다", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error for missing backtick name")
		}
	})

	t.Run("MissingMeans", func(t *testing.T) {
		if _, err := parseDefinitionKo(scan.Item{Text: "`x` 는 650 이다", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error for missing 뜻한다")
		}
	})

	t.Run("StructHeadKeywordVariants", func(t *testing.T) {
		children := []scan.Item{{Text: "`연소득` 을 가진다, 타입 Currency", No: 2}}
		bare, err := parseDefinitionKo(scan.Item{Text: "`신청인` 은 다음 뜻한다", No: 1, Children: children}, "t.md")
		if err != nil || bare.Kind != ast.StructDef || len(bare.Fields) != 1 {
			t.Fatalf("다음: got %+v err=%v", bare, err)
		}
		reul, err := parseDefinitionKo(scan.Item{Text: "`신청인` 은 다음를 뜻한다", No: 1, Children: children}, "t.md")
		if err != nil || reul.Kind != ast.StructDef || len(reul.Fields) != 1 {
			t.Fatalf("다음를: got %+v err=%v", reul, err)
		}
	})

	t.Run("StructHeadFieldError", func(t *testing.T) {
		it := scan.Item{Text: "`신청인` 은 다음을 뜻한다", No: 1, Children: []scan.Item{{Text: "bogus", No: 2}}}
		if _, err := parseDefinitionKo(it, "t.md"); err == nil {
			t.Fatal("expected error propagated from a malformed field")
		}
	})

	t.Run("ConstValueError", func(t *testing.T) {
		if _, err := parseDefinitionKo(scan.Item{Text: "`x` 는 650 을 뜻한다, `credit`.`Threshold`", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error propagated from a malformed const value tail")
		}
	})
}
