//ff:type feature=tangl type=model
//ff:what Name — the ten TANGL basic type name constants used in "as <Type>" clauses
package types

// Name is a TANGL basic type name, as written in Definitions ("as <Type>") and
// Cases/Provides ("is required as <Type>") clauses.
type Name = string

// The ten TANGL basic type names.
const (
	Text     Name = "Text"
	Integer  Name = "Integer"
	Number   Name = "Number"
	Boolean  Name = "Boolean"
	Email    Name = "Email"
	Date     Name = "Date"
	Time     Name = "Time"
	URL      Name = "URL"
	Currency Name = "Currency"
	Quantity Name = "Quantity"
)
