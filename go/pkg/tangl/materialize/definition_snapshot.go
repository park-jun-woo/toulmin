//ff:type feature=tangl type=model
//ff:what DefinitionSnapshot — portable loaded definition values shared by eval and gen
package materialize

// DefinitionSnapshot is the portable result shared by eval and gen.
type DefinitionSnapshot struct {
	Strings []string
}
