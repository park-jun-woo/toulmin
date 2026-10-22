//ff:type feature=tangl type=model
//ff:what See — a `see ... from ...` package reference declaration
package ast

// See declares a package reference: see `alias` from `package_path`.
type See struct {
	Alias   string `json:"alias"`
	Package string `json:"package"`
	Line    int    `json:"line"`
}
