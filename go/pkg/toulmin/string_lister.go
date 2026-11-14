//ff:type feature=engine type=interface
//ff:what StringLister — exposes a structured string list to comparison engines
package toulmin

// StringLister exposes a structured string list to comparison engines.
type StringLister interface {
	StringList() []string
}
