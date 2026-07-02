//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what multiError.Error — joins every collected violation onto its own line
package validate

import "strings"

// Error implements the error interface, joining every collected violation
// onto its own line.
func (m *multiError) Error() string {
	lines := make([]string, len(m.errs))
	for i, e := range m.errs {
		lines[i] = e.Error()
	}
	return strings.Join(lines, "\n")
}
