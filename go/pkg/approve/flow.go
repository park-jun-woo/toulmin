//ff:type feature=approve type=engine
//ff:what Flow: 다단계 승인 흐름
package approve

// Flow manages a multi-step approval workflow.
type Flow struct {
	name  string
	steps []*Step
}
