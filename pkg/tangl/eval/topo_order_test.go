//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what TestTopoOrder — tests dependency-first ordering with a shared, once-visited dep
package eval

import "testing"

func TestTopoOrder(t *testing.T) {
	names := []string{"a", "b", "c"}
	edges := map[string][]string{"a": {"c"}, "b": {"c"}, "c": nil}
	got := topoOrder(names, edges)
	want := []string{"c", "a", "b"}
	if len(got) != len(want) {
		t.Fatalf("want %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("want %v, got %v", want, got)
		}
	}
}
