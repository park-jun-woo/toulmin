//ff:type feature=engine type=model
//ff:what ruleIDAltSpec — spec sharing testSpec field layout but a different SpecName, to verify type identity drives the identifier (case b)
package toulmin

// ruleIDAltSpec shares the same field layout as testSpec (single Value string)
// but a different SpecName, to verify type identity drives the identifier (case b).
type ruleIDAltSpec struct{ Value string }
