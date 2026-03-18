//ff:func feature=engine type=util control=sequence
//ff:what FuncName — extracts short function name from function pointer
package toulmin

// FuncName returns the short name of a function from its pointer.
// e.g. "github.com/example/pkg.IsAdult" → "IsAdult"
func FuncName(fn func(any, any) (bool, any)) string {
	return shortName(funcID(fn))
}
