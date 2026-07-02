//ff:func feature=tangl type=engine control=selection
//ff:what Valid — checks that v satisfies the TANGL basic type named typeName
package types

// Valid reports whether v is a valid value of the TANGL basic type typeName.
// An unrecognized typeName always returns false.
func Valid(typeName string, v any) bool {
	switch typeName {
	case Text:
		return validText(v)
	case Integer:
		return validInteger(v)
	case Number:
		return validNumber(v)
	case Boolean:
		return validBoolean(v)
	case Email:
		return validEmail(v)
	case Date:
		return validDate(v)
	case Time:
		return validTime(v)
	case URL:
		return validURL(v)
	case Currency:
		return validCurrency(v)
	case Quantity:
		return validQuantity(v)
	default:
		return false
	}
}
