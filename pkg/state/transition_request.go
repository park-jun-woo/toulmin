//ff:type feature=state type=model
//ff:what TransitionRequest: 전이 요청 (from, to, event)
package state

// TransitionRequest represents a state transition request.
type TransitionRequest struct {
	From  string
	To    string
	Event string
}
