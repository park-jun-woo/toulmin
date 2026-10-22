//ff:type feature=approve type=adapter
//ff:what StepContextFunc: 단계별 ApprovalContext 생성 함수 타입
package approve

// StepContextFunc returns an ApprovalContext for the given step name.
type StepContextFunc func(stepName string) *ApprovalContext
