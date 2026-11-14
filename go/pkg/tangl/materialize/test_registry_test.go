//ff:type feature=tangl type=model
//ff:what testRegistry — test DefinitionLoaderRegistry that knows only the "JSON list" loader
package materialize

type testRegistry struct{ loader DefinitionLoader }
