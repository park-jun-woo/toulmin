//ff:func feature=feature type=engine control=sequence
//ff:what EvaluateTrace: 피처 판정 + 근거
package feature

import "fmt"

// EvaluateTrace returns the feature evaluation result with trace.
func (f *Flags) EvaluateTrace(name string, ctx *UserContext) (*FeatureResult, error) {
	g, ok := f.features[name]
	if !ok {
		return nil, fmt.Errorf("feature not registered: %s", name)
	}
	results, err := g.EvaluateTrace(name, ctx)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return &FeatureResult{Name: name, Enabled: false, Verdict: -1}, nil
	}
	return &FeatureResult{
		Name:    name,
		Enabled: results[0].Verdict > 0,
		Verdict: results[0].Verdict,
		Trace:   results[0].Trace,
	}, nil
}
