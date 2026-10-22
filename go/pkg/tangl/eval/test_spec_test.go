//ff:type feature=tangl type=model
//ff:what testSpec — test Spec with a name and a MinPrice field for SpecRef walking
package eval

// testSpec is a test toulmin.Spec: Name is its SpecName; MinPrice is a field a
// SpecRef can walk into.
type testSpec struct {
	Name     string
	MinPrice float64
}
