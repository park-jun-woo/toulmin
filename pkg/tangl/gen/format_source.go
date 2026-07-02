//ff:func feature=tangl type=codegen control=sequence
//ff:what formatSource — runs the assembled source through go/format.Source
package gen

import "go/format"

// formatSource runs the assembled source through go/format.Source, the
// final step Generate must pass before returning.
func formatSource(src string) (string, error) {
	out, err := format.Source([]byte(src))
	if err != nil {
		return "", err
	}
	return string(out), nil
}
