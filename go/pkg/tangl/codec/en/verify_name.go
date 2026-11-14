//ff:func feature=tangl type=parser control=sequence
//ff:what verifyName — check a scanned backtick name matches the AST slice value at a slot index
package en

import "fmt"

// verifyName reports an error when slice[i] does not equal name (or i is out of
// range), the drift guard for run/check body assignment (README decision 3).
func verifyName(slice []string, i int, name string) error {
	if i >= len(slice) || slice[i] != name {
		return fmt.Errorf("body name %q does not match AST", name)
	}
	return nil
}
