//ff:type feature=approve type=model
//ff:what ThresholdSpec — spec for amount threshold rules
package approve

// ThresholdSpec specifies an amount threshold.
type ThresholdSpec struct {
	Max float64 `yaml:"max"`
}
