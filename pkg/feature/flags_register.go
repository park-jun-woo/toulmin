//ff:func feature=feature type=engine control=sequence
//ff:what Register: 피처 graph 등록
package feature

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// Register adds a feature graph to the registry.
func (f *Flags) Register(name string, g *toulmin.Graph) {
	if _, exists := f.features[name]; !exists {
		f.order = append(f.order, name)
	}
	f.features[name] = g
}
