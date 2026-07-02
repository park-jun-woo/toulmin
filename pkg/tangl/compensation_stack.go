//ff:type feature=tangl type=model
//ff:what compensationStack — LIFO stack of compensation closures carried in a Context
package tangl

// compensationKey is the Context key under which the compensation stack is stored.
const compensationKey = "tangl:compensation"

// compensationStack holds the compensation closures armed by successful `do` actions
// during one top-level Run pass. Compensate runs fns in reverse (LIFO) order.
type compensationStack struct {
	fns []func() error
}
