//ff:func feature=tangl type=parser control=sequence
//ff:what renderSubjectKo — render the Subject declaration (`이 문서는 `x` 이다`)
package ko

// renderSubjectKo renders the tangl:주제 namespace back to its canonical form,
// "이 문서는 `<namespace>` 이다". Document.Subject is a Line-less plain string, so the
// value is quoted directly.
func renderSubjectKo(subject string) string {
	return "이 문서는 `" + subject + "` 이다"
}
