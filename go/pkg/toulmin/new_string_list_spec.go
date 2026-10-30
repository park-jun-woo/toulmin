//ff:func feature=engine type=util control=sequence
//ff:what NewStringListSpec — constructs a list Spec without retaining the caller's slice
package toulmin

// NewStringListSpec constructs a list Spec without retaining the caller's slice.
func NewStringListSpec(name string, items []string) StringListSpec {
	return StringListSpec{name: name, items: append([]string(nil), items...)}
}
