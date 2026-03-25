//ff:type feature=feature type=model
//ff:what AttributeSpec: HasAttribute rule의 spec 타입
package feature

// AttributeSpec carries key-value pair criteria for attribute checks.
type AttributeSpec struct {
	Key   string // attribute key
	Value any    // attribute value to match
}
