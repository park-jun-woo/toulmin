//ff:type feature=tangl type=util pattern=error-collection
//ff:what multiError — aggregates every Validate violation into a single error
package validate

// multiError aggregates every collected validation violation so Validate can
// return them as a single error while preserving each individual message.
type multiError struct {
	errs []error
}
