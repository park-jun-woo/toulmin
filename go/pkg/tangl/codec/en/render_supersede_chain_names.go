//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what renderSupersedeChainNames — render a supersede chain as backtick names joined by " > "
package en

// renderSupersedeChainNames renders the names of a supersede priority chain as
// "`A` > `B` > `C`" in chain order.
func renderSupersedeChainNames(names []string) string {
	text := ""
	for i, name := range names {
		if i > 0 {
			text += " > "
		}
		text += "`" + name + "`"
	}
	return text
}
