//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what resolveSpecs — resolves and validates specs for a rule from the registry
package toulmin

// resolveSpecs looks up specs by rule name and validates them if found.
func resolveSpecs(name string, specsRegistry map[string][]Spec) (Specs, error) {
	if specsRegistry == nil {
		return nil, nil
	}
	ss := specsRegistry[name]
	if len(ss) == 0 {
		return nil, nil
	}
	for _, s := range ss {
		if err := validateSpec(s); err != nil {
			return nil, err
		}
	}
	return ss, nil
}
