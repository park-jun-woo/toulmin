//ff:func feature=feature type=rule control=sequence
//ff:what HasAttribute: backing(AttributeBacking)으로 지정된 속성 키/값 쌍이 일치하는지 판정
package feature

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// HasAttribute checks if the user has the attribute key=value specified by backing.
func HasAttribute(claim any, ground any, backing toulmin.Backing) (bool, any) {
	ctx := ground.(*UserContext)
	ab := backing.(*AttributeBacking)
	return ctx.Attributes[ab.Key] == ab.Value, nil
}
