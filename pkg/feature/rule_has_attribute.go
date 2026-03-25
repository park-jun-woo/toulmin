//ff:func feature=feature type=rule control=sequence
//ff:what HasAttribute: spec(AttributeSpec)으로 지정된 속성 키/값 쌍이 일치하는지 판정
package feature

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// HasAttribute checks if the user has the attribute key=value specified by spec.
func HasAttribute(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	attributes, _ := ctx.Get("attributes")
	ab := specs[0].(*AttributeSpec)
	return attributes.(map[string]any)[ab.Key] == ab.Value, nil
}
