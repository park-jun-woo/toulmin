//ff:func feature=approve type=engine control=sequence
//ff:what NewFlow: 빈 Flow 생성
package approve

// NewFlow creates an empty Flow.
func NewFlow(name string) *Flow {
	return &Flow{name: name}
}
