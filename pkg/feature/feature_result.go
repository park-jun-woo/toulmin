//ff:type feature=feature type=model
//ff:what FeatureResult: 피처 판정 결과
package feature

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// FeatureResult holds the feature evaluation result with trace.
type FeatureResult struct {
	Name    string
	Enabled bool
	Verdict float64
	Trace   []toulmin.TraceEntry
}
