//ff:type feature=tangl type=codegen
//ff:what importSpec — one resolved Go import (optional alias + path)
package gen

// importSpec is one resolved Go import: Alias is the local identifier to
// bind (left blank when the import should use its path's own default
// name), Path is the full import path.
type importSpec struct {
	Alias string
	Path  string
}
