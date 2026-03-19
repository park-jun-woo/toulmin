//ff:func feature=approve type=engine control=sequence
//ff:what AddStep: 승인 단계 등록
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// AddStep registers an approval step.
func (f *Flow) AddStep(name string, g *toulmin.Graph) *Flow {
	f.steps = append(f.steps, &Step{Name: name, Graph: g})
	return f
}
