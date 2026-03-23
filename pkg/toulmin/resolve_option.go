//ff:func feature=engine type=engine control=sequence
//ff:what resolveOption — resolves EvalOption from variadic args with defaults
package toulmin

// resolveOption returns the first option or a zero-value default.
// Duration implies Trace.
func resolveOption(opts []EvalOption) EvalOption {
	if len(opts) == 0 {
		return EvalOption{}
	}
	opt := opts[0]
	if opt.Duration {
		opt.Trace = true
	}
	return opt
}
