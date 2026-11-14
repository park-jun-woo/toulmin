//ff:type feature=tangl type=model
//ff:what InternalKind — kind of a tangl:Internal trigger (event or tick)
package ast

// InternalKind classifies an Internal trigger.
type InternalKind int

const (
	// OnEvent is an `on <event>` trigger.
	OnEvent InternalKind = iota
	// EveryTick is an `every <interval> [until <case>]` trigger.
	EveryTick
)
