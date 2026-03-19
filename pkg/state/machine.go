//ff:type feature=state type=engine
//ff:what Machine: 전이 graph 등록 및 판정 실행
package state

// Machine manages state transition graphs.
type Machine struct {
	transitions map[string]*transition // key: "from:event"
	order       []string               // insertion order for Mermaid
}
