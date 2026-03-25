//ff:func feature=tangl type=validator control=sequence
//ff:what Validate — run all validation checks on a parsed TANGL File AST
package validate

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// Validate runs all checks on a parsed TANGL file and returns combined errors.
func Validate(f *parser.File) error {
	var errs []string
	errs = append(errs, checkDuplicates(f)...)
	errs = append(errs, checkGraphRefs(f)...)
	errs = append(errs, checkFuncRefs(f)...)
	errs = append(errs, checkAttackRefs(f)...)
	errs = append(errs, checkCycles(f)...)
	if len(errs) > 0 {
		return fmt.Errorf("validation errors:\n%s", strings.Join(errs, "\n"))
	}
	return nil
}
