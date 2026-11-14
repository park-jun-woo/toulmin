//ff:func feature=engine type=util control=iteration dimension=1
//ff:what TestBuildAttackerSet — tests buildAttackerSet for empty and populated edges branches
package toulmin

import "testing"

func TestBuildAttackerSet(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"Empty", func(t *testing.T) {
			set := buildAttackerSet(map[string][]string{})
			if len(set) != 0 {
				t.Fatalf("expected empty set, got %v", set)
			}
		}},
		{"Populated", func(t *testing.T) {
			edges := map[string][]string{
				"a": {"b", "c"},
				"b": {"c"},
			}
			set := buildAttackerSet(edges)
			if !set["b"] || !set["c"] || len(set) != 2 {
				t.Fatalf("expected set with b and c, got %v", set)
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
