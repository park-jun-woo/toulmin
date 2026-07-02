//ff:type feature=tangl type=model
//ff:what Logic — an "and"/"or" combination of sub-expressions
package ast

// Logic combines Terms with a boolean operator ("and" or "or").
type Logic struct {
	Op    string `json:"op"`
	Terms []Expr `json:"terms"`
}
