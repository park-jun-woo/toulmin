//ff:func feature=engine type=validator control=sequence
//ff:what validateBacking — validates a Backing instance (func fields + domain validation)
package toulmin

// validateBacking checks a backing for func fields and runs domain validation.
func validateBacking(b Backing) error {
	if err := validateBackingFields(b); err != nil {
		return err
	}
	return b.Validate()
}
