//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseDefinitionFieldsKo — collects fields, empty slice, propagates a child error
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestParseDefinitionFieldsKo(t *testing.T) {
	t.Run("Collects", func(t *testing.T) {
		items := []scan.Item{
			{Text: "`a` 를 가진다, 타입 Integer", No: 2},
			{Text: "`b` 를 가진다, 타입 Text", No: 3},
		}
		fields, err := parseDefinitionFieldsKo(items, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(fields) != 2 || fields[0].Name != "a" || fields[1].Type != "Text" {
			t.Fatalf("got %+v", fields)
		}
	})

	t.Run("Empty", func(t *testing.T) {
		fields, err := parseDefinitionFieldsKo(nil, "t.md")
		if err != nil || len(fields) != 0 {
			t.Fatalf("got %+v err=%v", fields, err)
		}
	})

	t.Run("ChildError", func(t *testing.T) {
		items := []scan.Item{{Text: "bogus", No: 2}}
		if _, err := parseDefinitionFieldsKo(items, "t.md"); err == nil {
			t.Fatal("expected error propagated from a malformed field")
		}
	})
}
