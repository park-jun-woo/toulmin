//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what snapshotLists — copy each materialized string-list Spec into a plain name → values map
package gen

import "github.com/park-jun-woo/toulmin/go/pkg/toulmin"

// snapshotLists returns the string values of every materialized list keyed by
// definition name, the embeddable snapshot generation writes into source.
func snapshotLists(materialized map[string]toulmin.StringListSpec) map[string][]string {
	lists := make(map[string][]string, len(materialized))
	for name, spec := range materialized {
		lists[name] = spec.StringList()
	}
	return lists
}
