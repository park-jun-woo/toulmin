//ff:func feature=feature type=rule control=sequence
//ff:what HasAttribute: backing([2]any)으로 지정된 속성 키/값 쌍이 일치하는지 판정
package feature

// HasAttribute checks if the user has the attribute key=value specified by backing.
// backing is [2]any{key string, value any}.
func HasAttribute(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*UserContext)
	pair := backing.([2]any)
	key := pair[0].(string)
	value := pair[1]
	return ctx.Attributes[key] == value, nil
}
