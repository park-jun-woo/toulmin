//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what packageName — converts the tangl:Subject value into a valid Go package name
package gen

// packageName converts the tangl:Subject value into a valid Go package
// name: any character outside [A-Za-z0-9_] becomes '_', and a leading
// digit is prefixed with '_' so the result is always a legal identifier.
func packageName(subject string) string {
	var b []rune
	for _, r := range subject {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			b = append(b, r)
		} else {
			b = append(b, '_')
		}
	}
	out := string(b)
	if out == "" {
		return "_"
	}
	if out[0] >= '0' && out[0] <= '9' {
		out = "_" + out
	}
	return out
}
