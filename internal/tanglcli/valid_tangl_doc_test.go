//ff:func feature=cli type=command control=sequence
//ff:what validTanglDoc — builds minimal-but-valid TANGL v0.3 markdown source declaring the given subject
package tanglcli

// validTanglDoc returns minimal-but-valid TANGL v0.3 markdown source whose
// tangl:Subject section declares the given subject.
func validTanglDoc(subject string) string {
	return "## tangl:Subject\n" +
		"- this document is `" + subject + "`\n" +
		"\n" +
		"## tangl:Cases\n" +
		"- in case of `c1`\n"
}
