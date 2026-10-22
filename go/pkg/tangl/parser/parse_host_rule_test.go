//ff:what TestParseHostRule — Rules host declaration owns using, with, and grounds
package parser

import "testing"

func TestParseHostRule(t *testing.T) {
	rule, err := parseInlineRuleItem(item{Text: "`tries exhausted` uses `ratchet`.`triesExhausted` with `max attempts` grounded on `attempts`", Line: 4}, "t.md")
	if err != nil {
		t.Fatal(err)
	}
	if rule.Using == nil || rule.Using.Alias != "ratchet" || rule.Using.Name != "triesExhausted" {
		t.Fatalf("using=%+v", rule.Using)
	}
	if len(rule.With) != 1 || rule.With[0] != "max attempts" || len(rule.Grounds) != 1 || rule.Grounds[0] != "attempts" {
		t.Fatalf("with=%v grounds=%v", rule.With, rule.Grounds)
	}
}
