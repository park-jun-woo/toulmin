//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what renderSupersedeChainNamesKo — render a supersede chain as backtick names joined by " > "
package ko

// renderSupersedeChainNamesKo renders the names of a supersede priority chain as
// "`A` > `B` > `C`" in chain order (the 순서로 배제한다 suffix is added by the caller).
func renderSupersedeChainNamesKo(names []string) string {
	text := ""
	for i, name := range names {
		if i > 0 {
			text += " > "
		}
		text += "`" + name + "`"
	}
	return text
}
