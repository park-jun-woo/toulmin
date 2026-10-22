//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what applyExecFixups — fill each deferred `run` target case name now that all names are known
package stream

// applyExecFixups resolves every deferred run-target reference: a `run` edge can
// point at a case declared later in the stream, so its target name is filled here
// from the completed node-name map.
func applyExecFixups(d *dec) {
	for _, f := range d.execCaseFixups {
		d.ast.Cases[f.caseIdx].Execs[f.slotIdx].Case = d.nodeName[f.tid]
	}
}
