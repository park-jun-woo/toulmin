//ff:type feature=tangl type=model
//ff:what DefKind — kind of a Definitions entry (const or struct)
package ast

// DefKind classifies a Definition as a constant or a struct.
type DefKind int

const (
	// ConstDef is a `` `x` means <literal> [as <Ref>] `` definition.
	ConstDef DefKind = iota
	// StructDef is a `` `x` means `` + nested `` has `f` as <Type> `` definition.
	StructDef
	// StringListDef is a syntax-directed JSON string-array Spec definition.
	StringListDef
	// FileSpecDef is a registry-loaded file-backed Spec definition.
	FileSpecDef
)
