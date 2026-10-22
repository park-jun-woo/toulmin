//ff:func feature=engine type=engine control=selection
//ff:what effectiveLevel — resolves a node's reporting Level (explicit set wins over role default)
package toulmin

// effectiveLevel returns the Level to record for a node: the explicitly set
// value when Level(...) was called (levelSetMap true), otherwise the role-based
// default. Used identically by the Evaluate (buildTraceEntry) and Run
// (recordTrace) trace paths so both report the same Level.
func (ec *evalContext) effectiveLevel(id, role string) Level {
	if ec.levelSetMap[id] {
		return ec.levelMap[id]
	}
	return defaultLevelForRole(role)
}
