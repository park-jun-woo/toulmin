//ff:type feature=engine type=model
//ff:what RunView — read-only snapshot of every node's final event after a full pass
package toulmin

// RunView is a read-only snapshot of the whole graph's final node events,
// shared with every handler during Run. It is immutable: a handler mutating
// ctx never changes the view another handler sees.
type RunView interface {
	All() []NodeEvent                  // every node's final event, registration order (Inactive included)
	Get(name string) (NodeEvent, bool) // node event by short name
	Attackers(name string) []NodeEvent // final events of nodes that attacked name
}
