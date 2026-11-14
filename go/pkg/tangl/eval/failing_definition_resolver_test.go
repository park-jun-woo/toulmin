//ff:type feature=tangl type=model
//ff:what failingDefinitionResolver — test DefinitionResolver returning a fixed error or a nil Spec
package eval

type failingDefinitionResolver struct {
	fakeResolver
	err error
}
