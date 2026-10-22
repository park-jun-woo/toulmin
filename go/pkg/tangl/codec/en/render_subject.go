//ff:func feature=tangl type=parser control=sequence
//ff:what renderSubject — render the Subject declaration (`this document is `x``)
package en

// renderSubject renders the tangl:Subject namespace back to its canonical form,
// "this document is `<namespace>`". Document.Subject is a Line-less plain string,
// so the value is quoted directly.
func renderSubject(subject string) string {
	return "this document is `" + subject + "`"
}
