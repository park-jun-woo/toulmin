//ff:type feature=tangl type=interface
//ff:what DefinitionLoader — reads and parses one file-backed definition
package materialize

// DefinitionLoader reads and parses one file-backed definition.
type DefinitionLoader interface {
	LoadDefinition(DefinitionLoadRequest) (DefinitionSnapshot, error)
}
