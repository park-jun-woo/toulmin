//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what renderPath — join backtick-quoted path steps with dots (`a`.`b`.`c`)
package en

// renderPath renders a data-path step list back to the backtick dot-chain
// surface form (`a`.`b`.`c`) — the inverse of the parser's parsePath. A single
// step renders as a lone backtick name, identical to a flat field.
func renderPath(steps []string) string {
	s := ""
	for i, step := range steps {
		if i > 0 {
			s += "."
		}
		s += "`" + step + "`"
	}
	return s
}
