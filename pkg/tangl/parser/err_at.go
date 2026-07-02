//ff:func feature=tangl type=parser control=sequence
//ff:what errAt — build a "path:line: message" formatted error
package parser

import "fmt"

// errAt formats an error consistently as "path:line: message".
func errAt(path string, line int, format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s:%d: %s", path, line, msg)
}
