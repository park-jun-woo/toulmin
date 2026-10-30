//ff:func feature=engine type=model control=sequence
//ff:what StringList — returns a copy of the list so callers cannot mutate the criterion
package toulmin

// StringList returns a copy so callers cannot mutate the criterion in place.
func (s StringListSpec) StringList() []string {
	return append([]string(nil), s.items...)
}
