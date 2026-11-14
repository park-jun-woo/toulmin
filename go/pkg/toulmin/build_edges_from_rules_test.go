//ff:func feature=engine type=util control=iteration dimension=1
//ff:what TestBuildEdgesFromRules — tests buildEdgesFromRules for no-rules, no-defeats, and multi-defeats branches
package toulmin

import (
	"reflect"
	"testing"
)

func TestBuildEdgesFromRules(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"NoRules", func(t *testing.T) {
			edges := map[string][]string{}
			buildEdgesFromRules(edges, nil)
			if len(edges) != 0 {
				t.Fatalf("expected empty edges, got %v", edges)
			}
		}},
		{"NoDefeats", func(t *testing.T) {
			edges := map[string][]string{}
			rules := []RuleMeta{
				{Name: "r1"},
			}
			buildEdgesFromRules(edges, rules)
			if len(edges) != 0 {
				t.Fatalf("expected empty edges, got %v", edges)
			}
		}},
		{"MultiDefeats", func(t *testing.T) {
			edges := map[string][]string{}
			rules := []RuleMeta{
				{Name: "r1", Defeats: []string{"a", "b"}},
				{Name: "r2", Defeats: []string{"a"}},
			}
			buildEdgesFromRules(edges, rules)
			want := map[string][]string{
				"a": {"r1", "r2"},
				"b": {"r1"},
			}
			if len(edges) != len(want) {
				t.Fatalf("expected %v, got %v", want, edges)
			}
			for k, v := range want {
				if !reflect.DeepEqual(edges[k], v) {
					t.Fatalf("edges[%q] = %v, want %v", k, edges[k], v)
				}
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
