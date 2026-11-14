//ff:type feature=tangl type=model
//ff:what testLoader — recording test DefinitionLoader (call count, last request, optional failure)
package materialize

type testLoader struct {
	calls   int
	fail    bool
	request DefinitionLoadRequest
}
