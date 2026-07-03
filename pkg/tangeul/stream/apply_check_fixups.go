//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what applyCheckFixups — fill each deferred `checking` target case name once all names are known
package stream

// applyCheckFixups resolves every deferred `checking` reference: verdict
// composition can target a case declared later, so the node's Checking field is
// filled here from the completed node-name map.
func applyCheckFixups(d *dec) {
	for _, f := range d.checkFixups {
		d.ast.Cases[f.caseIdx].Nodes[f.slotIdx].Checking = d.nodeName[f.tid]
	}
}
