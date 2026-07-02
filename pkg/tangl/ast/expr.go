//ff:type feature=tangl type=model
//ff:what Expr — condition expression node (Compare, Logic, Not, Either)
package ast

// Expr is a condition expression node. Concrete implementations are
// Compare, Logic, Not, and Either.
type Expr interface{ exprNode() }
