//ff:type feature=engine type=model
//ff:what StringListSpec — named immutable string-list judgment criterion
package toulmin

import "fmt"

// StringLister exposes a structured string list to comparison engines.
type StringLister interface {
	StringList() []string
}

// StringListSpec is the built-in Spec produced by syntax-directed list definitions.
type StringListSpec struct {
	name  string
	items []string
}

// NewStringListSpec constructs a list Spec without retaining the caller's slice.
func NewStringListSpec(name string, items []string) StringListSpec {
	return StringListSpec{name: name, items: append([]string(nil), items...)}
}

func (s StringListSpec) SpecName() string { return s.name }

func (s StringListSpec) Validate() error {
	if s.name == "" {
		return fmt.Errorf("string list Spec name is empty")
	}
	return nil
}

// StringList returns a copy so callers cannot mutate the criterion in place.
func (s StringListSpec) StringList() []string {
	return append([]string(nil), s.items...)
}
