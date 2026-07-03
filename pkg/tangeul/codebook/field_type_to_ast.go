//ff:func feature=tangl type=util control=selection
//ff:what FieldTypeToAST — map a FieldTypeCode back to its ast.Field.Type string
package codebook

import "fmt"

// FieldTypeToAST maps a FieldTypeCode back to its canonical ast.Field.Type
// string, the inverse of FieldTypeFromAST. An unrecognized code is a hard error
// — never a silent fallback.
func FieldTypeToAST(c FieldTypeCode) (string, error) {
	switch c {
	case FieldText:
		return "Text", nil
	case FieldInteger:
		return "Integer", nil
	case FieldNumber:
		return "Number", nil
	case FieldBoolean:
		return "Boolean", nil
	case FieldEmail:
		return "Email", nil
	case FieldDate:
		return "Date", nil
	case FieldTime:
		return "Time", nil
	case FieldURL:
		return "URL", nil
	case FieldCurrency:
		return "Currency", nil
	case FieldQuantity:
		return "Quantity", nil
	default:
		return "", fmt.Errorf("tangeul/codebook: unknown FieldTypeCode 0x%04X", uint16(c))
	}
}
