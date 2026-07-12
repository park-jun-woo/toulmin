//ff:type feature=tangl type=codegen
//ff:what binderInfo — codegen state of one in-scope exists binder
package gen

// binderInfo is the codegen-time state of one in-scope exists binder: the
// generated element and index variable names of the enclosing tanglAny
// element closure, and the source operand's Where notation. A binderEnv maps
// binder names to their binderInfo.
type binderInfo struct {
	Elem string
	Idx  string
	Src  reprCode
}
