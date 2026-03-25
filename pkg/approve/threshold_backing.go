//ff:type feature=approve type=model
//ff:what ThresholdBacking — backing for amount threshold rules
package approve

// ThresholdBacking specifies an amount threshold.
type ThresholdBacking struct {
	Max float64 `yaml:"max"`
}
