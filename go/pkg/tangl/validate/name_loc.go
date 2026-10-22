//ff:type feature=tangl type=model
//ff:what nameLoc — a declared name paired with its source line, for duplicate checks
package validate

// nameLoc pairs a declared name with its source line for duplicate-name
// and cycle-message reporting.
type nameLoc struct {
	Name string
	Line int
}
