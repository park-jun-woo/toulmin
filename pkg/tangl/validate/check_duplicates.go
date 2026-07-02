//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what checkDuplicates — reports every repeated name in locs, keyed to its first occurrence
package validate

// checkDuplicates walks locs in declaration order and reports every name seen
// more than once, referencing the line of its first declaration. kind labels
// the name space in the error message (e.g. "case name", "See alias").
func checkDuplicates(path, kind string, locs []nameLoc) []error {
	first := make(map[string]int)
	reported := make(map[string]bool)
	var errs []error
	for _, loc := range locs {
		line, ok := first[loc.Name]
		if !ok {
			first[loc.Name] = loc.Line
			continue
		}
		if reported[loc.Name] {
			continue
		}
		reported[loc.Name] = true
		errs = append(errs, errAt(path, loc.Line, "duplicate %s %q (first declared at line %d)", kind, loc.Name, line))
	}
	return errs
}
