//ff:func feature=tangl type=util control=sequence pattern=error-collection
//ff:what newMultiError — wraps collected violations into one error, or nil if none
package validate

// newMultiError returns nil when errs is empty, otherwise a single error
// whose message joins every collected violation (pattern=error-collection).
func newMultiError(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	return &multiError{errs: errs}
}
