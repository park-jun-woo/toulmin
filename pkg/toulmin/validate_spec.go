//ff:func feature=engine type=validator control=sequence
//ff:what validateSpec — validates a Spec instance (func fields + domain validation)
package toulmin

// validateSpec checks a spec for func fields and runs domain validation.
func validateSpec(s Spec) error {
	if err := validateSpecFields(s); err != nil {
		return err
	}
	return s.Validate()
}
