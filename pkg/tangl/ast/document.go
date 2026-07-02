//ff:type feature=tangl type=model
//ff:what Document — the full AST of a parsed TANGL v0.3 file
package ast

// Document is the parsed AST of a single TANGL v0.3 markdown file.
type Document struct {
	Path      string       `json:"path,omitempty"`
	Subject   string       `json:"subject"`
	Sees      []See        `json:"sees,omitempty"`
	Defs      []Definition `json:"defs,omitempty"`
	Rules     []InlineRule `json:"rules,omitempty"`
	Cases     []Case       `json:"cases"`
	Provides  []Endpoint   `json:"provides,omitempty"`
	Internals []Internal   `json:"internals,omitempty"`
}
