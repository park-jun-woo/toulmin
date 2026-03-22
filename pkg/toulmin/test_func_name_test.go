//ff:func feature=engine type=engine control=sequence
//ff:what TestFuncName — tests FuncName returns correct function name
package toulmin

import (
	"testing"
)

func TestFuncName(t *testing.T) {
	name := FuncName(WarrantA)
	if name != "WarrantA" {
		t.Errorf("expected 'WarrantA', got '%s'", name)
	}
}
