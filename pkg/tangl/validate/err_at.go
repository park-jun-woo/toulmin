//ff:func feature=tangl type=validator control=sequence
//ff:what errAt — build a "path:line: message" formatted error for validate/lint
package validate

import "fmt"

// errAt formats a validation violation consistently as "path:line: message".
func errAt(path string, line int, format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s:%d: %s", path, line, msg)
}
