//ff:func feature=state type=engine control=sequence
//ff:what NewMachine: 빈 Machine 생성
package state

// NewMachine creates an empty Machine.
func NewMachine() *Machine {
	return &Machine{
		transitions: make(map[string]*transition),
	}
}
