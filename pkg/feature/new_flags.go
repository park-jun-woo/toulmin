//ff:func feature=feature type=engine control=sequence
//ff:what NewFlags: 빈 Flags 레지스트리 생성
package feature

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// NewFlags creates an empty feature flag registry.
func NewFlags() *Flags {
	return &Flags{features: make(map[string]*toulmin.Graph)}
}
