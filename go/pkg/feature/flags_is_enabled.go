//ff:func feature=feature type=engine control=sequence
//ff:what IsEnabled: 피처 활성화 여부 판정
package feature

import "fmt"

// IsEnabled returns true if the feature is enabled for the given user context.
func (f *Flags) IsEnabled(name string, uctx *UserContext) (bool, error) {
	g, ok := f.features[name]
	if !ok {
		return false, fmt.Errorf("feature not registered: %s", name)
	}
	ctx := buildFeatureContext(name, uctx)
	results, err := g.Evaluate(ctx)
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, nil
	}
	return results[0].Verdict > 0, nil
}
