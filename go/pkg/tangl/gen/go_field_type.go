//ff:func feature=tangl type=codegen control=selection
//ff:what goFieldType — maps a Definitions field type name to its Go counterpart
package gen

// goFieldType maps a tangl:Definitions field type name to its Go
// counterpart: the tangl/types basic types map to plain Go types; any
// other name is assumed to be another Definitions struct and rendered as
// its PascalCase Go type name; an empty type defaults to "any".
func goFieldType(t string) string {
	switch t {
	case "Text", "Email", "Date", "Time", "URL":
		return "string"
	case "Integer":
		return "int"
	case "Number", "Currency", "Quantity":
		return "float64"
	case "Boolean":
		return "bool"
	case "":
		return "any"
	default:
		return goIdentExported(t)
	}
}
