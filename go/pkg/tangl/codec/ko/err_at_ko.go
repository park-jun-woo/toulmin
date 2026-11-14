//ff:func feature=tangl type=parser control=sequence
//ff:what errAtKo — build a "path:line: message" formatted error for KO parsers
package ko

import "fmt"

// errAtKo formats a KO parser error consistently as "path:line: message",
// mirroring pkg/tangl/parser.errAt so KO and EN diagnostics read the same.
func errAtKo(path string, line int, format string, args ...any) error {
	return fmt.Errorf("%s:%d: %s", path, line, fmt.Sprintf(format, args...))
}
