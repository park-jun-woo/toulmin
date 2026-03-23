//ff:func feature=engine type=engine control=sequence
//ff:what resolveBacking — resolves and validates backing for a rule from the registry
package toulmin

// resolveBacking looks up a backing by rule name and validates it if found.
func resolveBacking(name string, backings map[string]Backing) (Backing, error) {
	if backings == nil {
		return nil, nil
	}
	b := backings[name]
	if b == nil {
		return nil, nil
	}
	if err := validateBacking(b); err != nil {
		return nil, err
	}
	return b, nil
}
