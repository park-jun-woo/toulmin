//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseDefinitionFields — tests parseDefinitionFields for empty input, all-success accumulation, and first-error short-circuit
package parser

import "testing"

func TestParseDefinitionFields(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		fields, err := parseDefinitionFields(nil, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(fields) != 0 {
			t.Fatalf("expected no fields, got %+v", fields)
		}
	})

	t.Run("Success", func(t *testing.T) {
		children := []item{
			{Text: "has `name` as String", Line: 1},
			{Text: "has `age` as Int", Line: 2},
		}
		fields, err := parseDefinitionFields(children, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(fields) != 2 || fields[0].Name != "name" || fields[1].Name != "age" {
			t.Fatalf("expected fields [name, age], got %+v", fields)
		}
	})

	t.Run("ErrorStopsAtFirstFailure", func(t *testing.T) {
		children := []item{
			{Text: "has `name` as String", Line: 1},
			{Text: "bogus field", Line: 2},
			{Text: "has `age` as Int", Line: 3},
		}
		fields, err := parseDefinitionFields(children, "test.md")
		if err == nil {
			t.Fatal("expected an error for a malformed field item")
		}
		if fields != nil {
			t.Fatalf("expected nil fields on error, got %+v", fields)
		}
	})
}
