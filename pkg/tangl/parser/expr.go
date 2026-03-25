//ff:type feature=tangl type=model
//ff:what Expr — inline rule expression
package parser

// Expr represents a condition in an inline rule.
type Expr struct {
	Field    string // user, role, ip, amount
	OfSpec   bool   // true = spec field, false = ctx field
	Operator string // "is nil", "equals", "is greater than", etc.
	Value    any    // nil, "admin", 18, etc.
	And      *Expr  // chained and expression
	Or       *Expr  // chained or expression
}
