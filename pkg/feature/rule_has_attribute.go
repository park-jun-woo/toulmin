//ff:func feature=feature type=rule control=sequence
//ff:what HasAttribute: backing(AttributeBacking)으로 지정된 속성 키/값 쌍이 일치하는지 판정
package feature

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// HasAttribute checks if the user has the attribute key=value specified by backing.
func HasAttribute(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	attributes, _ := ctx.Get("attributes")
	ab := backing.(*AttributeBacking)
	return attributes.(map[string]any)[ab.Key] == ab.Value, nil
}
