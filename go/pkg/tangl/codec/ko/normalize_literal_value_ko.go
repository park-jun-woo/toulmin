//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what normalizeLiteralValueKo — strip literal quotes from a KO comparison value (each item for "is in")
package ko

import "strings"

// normalizeLiteralValueKo strips one pair of surrounding double quotes from a
// literal comparison value. For "is in" the value is a comma list: every item is
// trimmed and unquoted, then re-joined with ", ".
func normalizeLiteralValueKo(op, value string) string {
	if op != "is in" {
		return stripLiteralQuotesKo(value)
	}
	items := strings.Split(value, ",")
	for i := range items {
		items[i] = stripLiteralQuotesKo(strings.TrimSpace(items[i]))
	}
	return strings.Join(items, ", ")
}
