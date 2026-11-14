//ff:func feature=tangl type=engine control=selection
//ff:what logicClosure — selects the and/or closure for a Logic operator
package eval

import "fmt"

// logicClosure returns the closure for the given Logic operator, or an error
// for an unknown operator.
func logicClosure(op string, terms []evalFunc) (evalFunc, error) {
	switch op {
	case "and":
		return andClosure(terms), nil
	case "or":
		return orClosure(terms), nil
	default:
		return nil, fmt.Errorf("eval: unknown logic operator %q", op)
	}
}
