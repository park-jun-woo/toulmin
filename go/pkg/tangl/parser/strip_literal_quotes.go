//ff:func feature=tangl type=parser control=sequence
//ff:what stripLiteralQuotes — remove one pair of surrounding double quotes from a literal
package parser

// stripLiteralQuotes returns value without its outer double quotes when it
// both starts and ends with '"' (length >= 2); otherwise value unchanged.
func stripLiteralQuotes(value string) string {
	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
		return value[1 : len(value)-1]
	}
	return value
}
