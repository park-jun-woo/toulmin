//ff:type feature=feature type=model
//ff:what AttributeBacking: HasAttribute rule의 backing 타입
package feature

// AttributeBacking carries key-value pair criteria for attribute checks.
type AttributeBacking struct {
	Key   string // attribute key
	Value any    // attribute value to match
}
