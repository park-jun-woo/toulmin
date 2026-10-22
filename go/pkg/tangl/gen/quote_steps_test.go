//ff:func feature=tangl type=codegen control=sequence
//ff:what TestQuoteSteps — tests Go string-literal list rendering of path steps
package gen

import "testing"

func TestQuoteSteps(t *testing.T) {
	if got := quoteSteps(nil); got != "" {
		t.Errorf("no steps: got %q", got)
	}
	if got := quoteSteps([]string{"author", "trustScore"}); got != `"author", "trustScore"` {
		t.Errorf("two steps: got %q", got)
	}
}
