//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what normalizeLiteralValue — strip literal quotes from a comparison value (each item for "is in")
package parser

import "strings"

// normalizeLiteralValue strips one pair of surrounding double quotes from a
// literal comparison value. For "is in" the value is a comma list: every item
// is trimmed and unquoted, then re-joined with ", ".
func normalizeLiteralValue(op, value string) string {
	if op != "is in" {
		return stripLiteralQuotes(value)
	}
	items := strings.Split(value, ",")
	for i := range items {
		items[i] = stripLiteralQuotes(strings.TrimSpace(items[i]))
	}
	return strings.Join(items, ", ")
}
