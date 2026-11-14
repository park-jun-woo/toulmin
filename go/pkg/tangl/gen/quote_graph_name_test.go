//ff:func feature=tangl type=codegen control=sequence dimension=1
//ff:what TestQuoteGraphName — tests quoteGraphName formats a quoted subject.case string
package gen

import "testing"

func TestQuoteGraphName(t *testing.T) {
	got := quoteGraphName("orders", "caseA")
	want := `"orders.caseA"`
	if got != want {
		t.Errorf("quoteGraphName() = %q, want %q", got, want)
	}
}
