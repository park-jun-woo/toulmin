//ff:type feature=feature type=engine
//ff:what Flags: 피처 플래그 레지스트리
package feature

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// Flags manages feature flag graphs.
type Flags struct {
	features map[string]*toulmin.Graph
	order    []string
}
