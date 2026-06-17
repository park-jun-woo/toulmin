//ff:type feature=engine type=model
//ff:what NodeEvent — payload passed to a node event handler
package toulmin

// NodeEvent is the payload delivered to a node's event handler.
type NodeEvent struct {
	Name     string        // node short name
	Role     string        // "rule" | "counter" | "except"
	Type     NodeEventType // Inactive | Active | Defeated
	Verdict  float64       // final verdict (0.0 when Inactive)
	Evidence any           // evidence produced by the rule function
}
