//ff:func feature=tangl type=util control=selection
//ff:what FieldTypeToAST — map a FieldTypeCode back to its ast.Field.Type string
package codebook

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/toulmin"
)

// FieldTypeToAST maps a FieldTypeCode back to its canonical ast.Field.Type
// string, the inverse of FieldTypeFromAST. An unrecognized code is a hard error
// — never a silent fallback.
func FieldTypeToAST(c toulmin.FieldTypeCode) (string, error) {
	switch c {
	case toulmin.FieldText:
		return "Text", nil
	case toulmin.FieldInteger:
		return "Integer", nil
	case toulmin.FieldNumber:
		return "Number", nil
	case toulmin.FieldBoolean:
		return "Boolean", nil
	case toulmin.FieldEmail:
		return "Email", nil
	case toulmin.FieldDate:
		return "Date", nil
	case toulmin.FieldTime:
		return "Time", nil
	case toulmin.FieldURL:
		return "URL", nil
	case toulmin.FieldCurrency:
		return "Currency", nil
	case toulmin.FieldQuantity:
		return "Quantity", nil
	default:
		return "", fmt.Errorf("tangeul/codebook: unknown toulmin.FieldTypeCode 0x%04X", uint16(c))
	}
}
