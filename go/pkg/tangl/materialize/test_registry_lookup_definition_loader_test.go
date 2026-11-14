//ff:func feature=tangl type=engine control=sequence
//ff:what LookupDefinitionLoader — testRegistry resolves only the "JSON list" type name
package materialize

func (r testRegistry) LookupDefinitionLoader(name string) (DefinitionLoader, bool) {
	return r.loader, name == "JSON list"
}
