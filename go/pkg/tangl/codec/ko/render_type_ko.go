//ff:func feature=tangl type=parser control=sequence
//ff:what renderTypeKo — render a '타입' clause token (bare built-in or `alias`.`Name` ref)
package ko

import "strings"

// renderTypeKo renders a TYPE_REF token (as stored by parseTypeRefKo) back to its
// KO surface form: a qualified reference "alias.Name" becomes `alias`.`Name`, while
// a bare token (a built-in type such as Currency, or an unqualified name) is emitted
// verbatim. Both re-parse to the same stored string.
func renderTypeKo(typ string) string {
	if i := strings.IndexByte(typ, '.'); i >= 0 {
		return "`" + typ[:i] + "`.`" + typ[i+1:] + "`"
	}
	return typ
}
