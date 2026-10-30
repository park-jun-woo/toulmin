//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseInlineRuleItemHostErrors — host rule declarations reject bad refs, qualified weights and trailing text
package parser

import (
	"strings"
	"testing"
)

func TestParseInlineRuleItemHostErrors(t *testing.T) {
	for _, c := range []struct{ text, want string }{
		{"`r` uses nothing", "expected function reference after 'using'"},
		{"`r` uses `ratchet`.`f` qualified 0.5", "qualified belongs to a case binding"},
		{"`r` uses `ratchet`.`f` and more", "unexpected trailing text"},
	} {
		_, err := parseInlineRuleItem(item{Text: c.text, Line: 6}, "t.md")
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%q: err=%v, want %q", c.text, err, c.want)
		}
	}
}
