//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what appendUniqueStrings — append strings once while preserving first occurrence order
package gen

func appendUniqueStrings(dst []string, values ...string) []string {
	seen := make(map[string]bool, len(dst)+len(values))
	for _, value := range dst {
		seen[value] = true
	}
	for _, value := range values {
		if !seen[value] {
			seen[value] = true
			dst = append(dst, value)
		}
	}
	return dst
}
