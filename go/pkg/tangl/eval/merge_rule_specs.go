//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what mergeRuleSpecs — merge rule-owned Specs with caller-supplied Specs, rejecting same-name conflicts
package eval

import (
	"fmt"
	"reflect"

	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// mergeRuleSpecs returns owned followed by every supplied Spec whose name is not
// already present. A supplied Spec sharing a name with an existing one must be
// deeply equal to it; otherwise the merge fails.
func mergeRuleSpecs(owned, supplied toulmin.Specs) (toulmin.Specs, error) {
	merged := append(toulmin.Specs(nil), owned...)
	for _, candidate := range supplied {
		existing := merged.Find(candidate.SpecName())
		if existing == nil {
			merged = append(merged, candidate)
			continue
		}
		if !reflect.DeepEqual(existing, candidate) {
			return nil, fmt.Errorf("conflicting Spec %q supplied to rule", candidate.SpecName())
		}
	}
	return merged, nil
}
