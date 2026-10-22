//ff:func feature=tangl type=parser control=sequence
//ff:what renderType — render an "as Type" token (bare built-in or `alias`.`Name` ref)
package en

import "strings"

// renderType renders a TYPE_REF token (as stored by the parser's parseTypeRef) back
// to its surface form: a qualified reference "alias.Name" becomes `alias`.`Name`,
// while a bare token (a built-in type such as Currency, or an unqualified name)
// is emitted verbatim. Both re-parse to the same stored string.
func renderType(typ string) string {
	if i := strings.IndexByte(typ, '.'); i >= 0 {
		return "`" + typ[:i] + "`.`" + typ[i+1:] + "`"
	}
	return typ
}
