//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what renderGroundsKo — render a Korean grounded-on adnominal name list
package ko

func renderGroundsKo(names []string) string {
	text := ""
	for i, name := range names {
		if i > 0 {
			text += " 그리고 "
		}
		text += "`" + name + "`"
	}
	if len(names) == 0 {
		return ""
	}
	return text + " " + josaFor(names[len(names)-1], "을", "를") + " 근거로 하는"
}
