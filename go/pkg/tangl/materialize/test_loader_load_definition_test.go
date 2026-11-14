//ff:func feature=tangl type=engine control=sequence
//ff:what LoadDefinition — testLoader records the request and returns a fixed snapshot or an error
package materialize

import "errors"

func (l *testLoader) LoadDefinition(r DefinitionLoadRequest) (DefinitionSnapshot, error) {
	l.calls++
	l.request = r
	if l.fail {
		return DefinitionSnapshot{}, errors.New("broken")
	}
	return DefinitionSnapshot{Strings: []string{"a,b", "c"}}, nil
}
