//ff:func feature=tangl type=codegen control=sequence dimension=1
//ff:what TestOnceKey — tests onceKey formats the composite once-guard key
package gen

import "testing"

func TestOnceKey(t *testing.T) {
	got := onceKey("orders", "caseA", "nodeB", 3)
	want := "once:orders.caseA.nodeB#3"
	if got != want {
		t.Errorf("onceKey() = %q, want %q", got, want)
	}
}
