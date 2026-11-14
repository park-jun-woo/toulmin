//ff:type feature=tangl type=interface
//ff:what DefinitionLoaderRegistry — resolves a surface loader type name to a DefinitionLoader
package materialize

// DefinitionLoaderRegistry resolves a surface loader type name.
type DefinitionLoaderRegistry interface {
	LookupDefinitionLoader(typeName string) (DefinitionLoader, bool)
}
