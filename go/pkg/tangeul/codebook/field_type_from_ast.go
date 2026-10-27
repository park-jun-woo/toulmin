//ff:func feature=tangl type=util control=selection
//ff:what FieldTypeFromAST — map an ast.Field.Type string to its FieldTypeCode
package codebook

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// FieldTypeFromAST maps a canonical ast.Field.Type string (one of the ten
// TANGEUL basic type names) to its FieldTypeCode. An unrecognized string is a
// hard error — never a silent fallback.
func FieldTypeFromAST(name string) (toulmin.FieldTypeCode, error) {
	switch name {
	case "Text":
		return toulmin.FieldText, nil
	case "Integer":
		return toulmin.FieldInteger, nil
	case "Number":
		return toulmin.FieldNumber, nil
	case "Boolean":
		return toulmin.FieldBoolean, nil
	case "Email":
		return toulmin.FieldEmail, nil
	case "Date":
		return toulmin.FieldDate, nil
	case "Time":
		return toulmin.FieldTime, nil
	case "URL":
		return toulmin.FieldURL, nil
	case "Currency":
		return toulmin.FieldCurrency, nil
	case "Quantity":
		return toulmin.FieldQuantity, nil
	default:
		return 0, fmt.Errorf("tangeul/codebook: unknown field type %q", name)
	}
}
