//ff:func feature=tangl type=util control=selection
//ff:what FieldTypeFromAST — map an ast.Field.Type string to its FieldTypeCode
package codebook

import "fmt"

// FieldTypeFromAST maps a canonical ast.Field.Type string (one of the ten
// TANGEUL basic type names) to its FieldTypeCode. An unrecognized string is a
// hard error — never a silent fallback.
func FieldTypeFromAST(name string) (FieldTypeCode, error) {
	switch name {
	case "Text":
		return FieldText, nil
	case "Integer":
		return FieldInteger, nil
	case "Number":
		return FieldNumber, nil
	case "Boolean":
		return FieldBoolean, nil
	case "Email":
		return FieldEmail, nil
	case "Date":
		return FieldDate, nil
	case "Time":
		return FieldTime, nil
	case "URL":
		return FieldURL, nil
	case "Currency":
		return FieldCurrency, nil
	case "Quantity":
		return FieldQuantity, nil
	default:
		return 0, fmt.Errorf("tangeul/codebook: unknown field type %q", name)
	}
}
