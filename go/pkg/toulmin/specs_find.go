//ff:func feature=engine type=model control=iteration dimension=1
//ff:what Find — looks up a Spec by SpecName
package toulmin

// Find returns the first Spec matching the given name, or nil.
func (s Specs) Find(name string) Spec {
	for _, sp := range s {
		if sp.SpecName() == name {
			return sp
		}
	}
	return nil
}
