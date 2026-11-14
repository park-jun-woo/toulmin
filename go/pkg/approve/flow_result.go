//ff:type feature=approve type=model
//ff:what FlowResult: 다단계 승인 판정 결과
package approve

// FlowResult holds the result of a multi-step approval flow.
type FlowResult struct {
	Approved bool
	Steps    []StepResult
}
