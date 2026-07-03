//ff:func feature=tangl type=parser control=sequence
//ff:what bumpCounter — advance the run or check counter for the just-assigned slot kind
package en

// bumpCounter returns runN and checkN with the counter for slotKind (1 Run, 2
// Check) advanced by one, keeping the per-kind slot indexing in assignBodySlots.
func bumpCounter(slotKind, runN, checkN int) (int, int) {
	if slotKind == 1 {
		return runN + 1, checkN
	}
	return runN, checkN + 1
}
