//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what moduleLine — extracts the module path from a go.mod file's "module" directive
package gen

import "strings"

// moduleLine scans go.mod bytes for the first "module <path>" directive and
// returns the declared module path.
func moduleLine(b []byte) (string, bool) {
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), true
		}
	}
	return "", false
}
