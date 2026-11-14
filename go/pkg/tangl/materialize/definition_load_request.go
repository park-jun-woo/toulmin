//ff:type feature=tangl type=model
//ff:what DefinitionLoadRequest — one file-backed definition load (name, loader type, source and resolved path)
package materialize

// DefinitionLoadRequest describes one file-backed definition load.
type DefinitionLoadRequest struct {
	DefinitionName string
	TypeName       string
	SourcePath     string
	ResolvedPath   string
}
